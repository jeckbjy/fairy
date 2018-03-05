package fairy

import (
	"fairy/container/inlist"
	"fairy/util"
	"sync"
	"sync/atomic"
	"time"
)

const (
	TIME_INTERVAL = 1             // 默认时间间隔
	WHEEL_NUM     = 3             // 初始wheel个数(2^30)，可以扩展
	SLOT_POW      = 10            // 2^SLOT_POW
	SLOT_MAX      = 1 << SLOT_POW // 个数
)

var gAsyncTimerEngine *TimerEngine
var gTimerEngine *TimerEngine

func GetAsyncTimerEngine() *TimerEngine {
	util.Once(gAsyncTimerEngine, func() {
		gAsyncTimerEngine = NewTimerEngine(nil)
		gAsyncTimerEngine.Start()
	})

	return gAsyncTimerEngine
}

func GetTimerEngine() *TimerEngine {
	util.OnceEx(gTimerEngine, func() {
		gTimerEngine = NewTimerEngine(GetExecutor())
		gTimerEngine.Start()
	})

	return gTimerEngine
}

func NewTimerEngine(exec *Executor) *TimerEngine {
	engine := &TimerEngine{}
	engine.Create(exec)
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
// TODO:AddTimer and DelTimer is not thread safe!!!
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

func (self *TimerEngine) Create(exec *Executor) {
	self.SetInterval(TIME_INTERVAL)
	self.executor = exec
	self.timestamp = util.Now()
	self.count = 0
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
	self.interval = int64(interval)
}

func (self *TimerEngine) AddTimer(timer *Timer) {
	if timer.Timestamp > self.timestamp {
		self.mutex.Lock()
		atomic.AddInt32(&self.count, 1)
		self.push(timer)
		self.mutex.Unlock()
	}
}

func (self *TimerEngine) DelTimer(timer *Timer) {
	self.mutex.Lock()
	if timer.IsRunning() {
		atomic.AddInt32(&self.count, -1)
		timer.List().Remove(timer)
	}
	self.mutex.Unlock()
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
	}
}

func (self *TimerEngine) OnExit() {
	self.Stop()
}

func (self *TimerEngine) Update() {
	oldTime := self.timestamp
	newTime := util.Now()
	self.timestamp = newTime
	if newTime <= oldTime {
		// 可能由于重新调时间导致
		self.build(oldTime, newTime)
	} else {
		self.tick(oldTime, newTime)
	}

	// check timer change
}

func (self *TimerEngine) build(oldTime int64, newTime int64) {
	// reset all timer
	self.mutex.Lock()
	pendings := inlist.List{}
	for _, wheel := range self.wheels {
		wheel.index = 0
		for _, slot := range wheel.slots {
			pendings.MoveBackList(slot)
		}
	}

	//
	for iter := pendings.Front(); iter != nil; {
		timer := iter.(*Timer)
		iter = inlist.Next(iter)
		timer.reset(oldTime, newTime)

		//
		if timer.Timestamp >= newTime {
			pendings.Remove(timer)
			self.push(timer)
		}
	}

	if self.executor != nil && pendings.Len() > 0 {
		self.pendings.MoveBackList(&pendings)
		self.executor.Dispatch(NewTimerEvent(self))
	}

	self.mutex.Unlock()

	// async call
	for iter := pendings.Front(); iter != nil; {
		timer := iter.(*Timer)
		iter = inlist.Next(iter)
		timer.Invoke()
		atomic.AddInt32(&self.count, -1)
	}
}

func (self *TimerEngine) tick(oldTime int64, newTime int64) {
	self.mutex.Lock()
	pendings := inlist.List{}
	// tick all timer
	for ticks := newTime - oldTime; ticks > 0; ticks-- {
		wheel := self.wheels[0]
		// process timer list
		pendings.MoveBackList(wheel.Current())
		// step tick
		self.cascade()
	}

	if self.executor != nil && pendings.Len() > 0 {
		self.pendings.MoveBackList(&pendings)
		self.executor.Dispatch(NewTimerEvent(self))
	}

	self.mutex.Unlock()

	// async call
	for iter := pendings.Front(); iter != nil; {
		timer := iter.(*Timer)
		iter = inlist.Next(iter)
		timer.Invoke()
		atomic.AddInt32(&self.count, -1)
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
	for iter := pendings.Front(); iter != nil; iter = inlist.Next(iter) {
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
