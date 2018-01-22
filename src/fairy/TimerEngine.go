package fairy

import (
	"fairy/container/inlist"
	"fairy/util"
	"sync"
	"sync/atomic"
	"time"
)

const (
	WHEEL_NUM = 3             // 初始wheel个数(2^30)，可以扩展
	SLOT_POW  = 10            // 2^SLOT_POW
	SLOT_MAX  = 1 << SLOT_POW // 个数
)

var gTimerEngine *TimerEngine

func GetTimerEngine() *TimerEngine {
	if gTimerEngine == nil {
		gTimerEngine = NewTimerEngine()
	}

	return gTimerEngine
}

func NewTimerEngine() *TimerEngine {
	engine := &TimerEngine{}
	engine.Create()
	return engine
}

//////////////////////////////////////////////////////
// TimerWheel
//////////////////////////////////////////////////////
type Wheel struct {
	slots   []*inlist.List // 桶
	index   int            // 循环索引
	timeOff uint
	timeMax uint64
}

func (self *Wheel) Create(index int) {
	self.index = 0
	self.timeOff = uint(index * SLOT_POW)
	self.timeMax = uint64(1) << uint((index+1)*SLOT_POW)
	for i := 0; i < SLOT_MAX; i++ {
		self.slots = append(self.slots, inlist.New())
	}
}

func (self *Wheel) Current() *inlist.List {
	return self.slots[self.index]
}

func (self *Wheel) Step() bool {
	self.index++
	if self.index >= len(self.slots) {
		self.index = 0
		return true
	}

	return false
}

func (self *Wheel) Push(timer *Timer, delta uint64) {
	index := (int(delta>>self.timeOff) + self.index) % len(self.slots)
	self.slots[index].PushBack(timer)
}

//////////////////////////////////////////////////////
// TimerEngine
// https://www.cnblogs.com/zhanghairong/p/3757656.html
//////////////////////////////////////////////////////
type TimerEngine struct {
	timestamp int64
	interval  int64
	count     int32
	wheels    []*Wheel
	executor  *Executor
	pendings  *inlist.List
	mutex     *sync.Mutex
	stopped   bool
	timeMax   uint64
}

func (self *TimerEngine) Create() {
	self.SetInterval(1)
	self.timestamp = 0
	self.count = 0
	self.executor = GetExecutor()
	self.pendings = inlist.New()
	self.stopped = true
	self.mutex = &sync.Mutex{}

	// create wheel
	for i := 0; i < WHEEL_NUM; i++ {
		self.newWheel()
	}
}

func (self *TimerEngine) newWheel() {
	wheel := &Wheel{}
	wheel.Create(len(self.wheels))
	self.timeMax = wheel.timeMax
	self.wheels = append(self.wheels, wheel)
}

func (self *TimerEngine) SetExecutor(exec *Executor) {
	self.executor = exec
}

func (self *TimerEngine) SetInterval(interval int) {
	self.interval = int64(interval) * int64(time.Millisecond)
}

func (self *TimerEngine) AddTimer(timer *Timer) {
	if timer.Timestamp > self.timestamp {
		atomic.AddInt32(&self.count, 1)
		self.push(timer)
	}
}

func (self *TimerEngine) DelTimer(timer *Timer) {
	atomic.AddInt32(&self.count, -1)
	timer.List().Remove(timer)
}

func (self *TimerEngine) Run() {
	for !self.stopped {
		self.Update()
		time.Sleep(time.Duration(self.interval) * time.Millisecond)
	}
}

func (self *TimerEngine) Start() {
	if self.stopped {
		RegisterExit(self)
		self.stopped = false
		go self.Run()
	}
}

func (self *TimerEngine) Stop() {
	if !self.stopped {
		self.stopped = true
		// notify
	}
}

func (self *TimerEngine) OnExit() {
	self.Stop()
}

func (self *TimerEngine) Update() {
	now := util.Now()
	self.Tick(now)
}

func (self *TimerEngine) Tick(now int64) {
	oldTimestamp := self.timestamp
	self.timestamp = now

	pendings := inlist.List{}
	if now < self.timestamp {
		// reset all timer
		timers := inlist.List{}
		for _, wheel := range self.wheels {
			for _, slot := range wheel.slots {
				timers.MoveBackList(slot)
			}
		}

		// reset
		for iter := timers.Front(); iter != nil; iter = inlist.Next(iter) {
			timer := iter.(*Timer)
			iter = inlist.Next(iter)
			timers.Remove(timer)

			// reset timer
			if timer.Delay > 0 {
				elapse := timer.Timestamp - oldTimestamp
				left := int64(timer.Delay) - elapse
				if left > 0 {
					timer.Delay = int(left)
					timer.Timestamp = now + left
				}
			}

			if timer.Timestamp > now {
				// push back
				self.push(timer)
			} else {
				// check timestamp
				if self.executor == nil || timer.Async {
					timer.Invoke()
					atomic.AddInt32(&self.count, -1)
				} else {
					pendings.PushBack(timer)
				}
			}
		}

	} else {
		// tick all timer
		ticks := now - oldTimestamp
		for ticks > 0 {
			wheel := self.wheels[0]
			// process timer list
			timers := wheel.Current()
			for iter := timers.Front(); iter != nil; {
				timer := iter.(*Timer)
				iter = inlist.Next(iter)
				timers.Remove(timer)

				if self.executor == nil || timer.Async {
					timer.Invoke()
					atomic.AddInt32(&self.count, -1)
				} else {
					pendings.PushBack(timer)
				}
			}

			// step tick
			self.cascade()
		}
	}

	// pendings for executor invoke
	if pendings.Len() > 0 {
		self.mutex.Lock()
		self.pendings.MoveBackList(&pendings)
		self.executor.Dispatch(NewTimerEvent(self))
		self.mutex.Unlock()
	}
}

func (self *TimerEngine) cascade() {
	for i := 0; i < len(self.wheels); i++ {
		if !self.wheels[i].Step() {
			break
		}

		// calc next
		if i+1 == len(self.wheels) {
			self.newWheel()
		} else {
			slots := self.wheels[i].Current()
			for iter := slots.Front(); iter != nil; {
				timer := iter.(*Timer)
				iter = inlist.Next(iter)
				slots.Remove(timer)

				self.push(timer)
			}
		}
	}
}

func (self *TimerEngine) Invoke() {
	// invoke all pending timers
	pendings := inlist.List{}
	self.mutex.Lock()
	pendings = *self.pendings
	self.pendings.Init()
	self.mutex.Unlock()

	// invoke
	for iter := pendings.Front(); iter != nil; inlist.Next(iter) {
		timer := iter.(*Timer)
		timer.Invoke()
		atomic.AddInt32(&self.count, -1)
	}
}

func (self *TimerEngine) push(timer *Timer) {
	if timer.Timestamp <= self.timestamp {
		atomic.AddInt32(&self.count, -1)
		return
	}

	delta := uint64(timer.Timestamp - self.timestamp)
	// check dynamic create wheel
	for delta > self.timeMax {
		self.newWheel()
	}

	for i := 0; i < len(self.wheels); i++ {
		wheel := self.wheels[i]
		if delta < wheel.timeMax {
			wheel.Push(timer, delta)
			break
		}
	}
}
