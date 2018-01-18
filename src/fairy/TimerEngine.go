package fairy

import (
	"fairy/container/inlist"
	"fairy/util"
	"sync"
	"time"
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
	engine.executor = GetExecutor()
	return engine
}

type Wheel struct {
	slots   []*inlist.List
	index   int
	offset  uint
	timeMax int64
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

type TimerEngine struct {
	start     int64
	timestamp int64
	interval  int64
	count     int
	wheels    []*Wheel
	executor  *Executor
	pendings  *inlist.List
	mutex     sync.Mutex
	stopped   bool
}

func (self *TimerEngine) SetExecutor(exec *Executor) {
	self.executor = exec
}

func (self *TimerEngine) SetInterval(interval int) {
	self.interval = int64(interval) * int64(time.Millisecond)
}

func (self *TimerEngine) AddTimer(timer *Timer) {
	self.count++
	self.push(timer)
}

func (self *TimerEngine) DelTimer(timer *Timer) {
	self.count--
	// remove by ?
}

func (self *TimerEngine) Start() {
	if self.stopped {
		self.stopped = false
		go func() {
			self.Update()
			time.Sleep(time.Duration(self.interval) * time.Millisecond)
		}()
	}
}

func (self *TimerEngine) Stop() {
	if !self.stopped {
		self.stopped = true
		// notify
	}
}

func (self *TimerEngine) Update() {
	now := util.Now()
	self.Tick(now)
}

func (self *TimerEngine) Tick(now int64) {
	pendings := inlist.List{}
	if now < self.timestamp {
		// reset all timer
		// timers := &List{}
		// for i := 0; i < len(self.wheels); i++ {
		// 	wheel := self.wheels[i]
		// 	// all slots

		// }
	} else {
		// tick all timer
		ticks := now - self.timestamp
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
				} else {
					pendings.PushBack(timer)
				}
			}

			// step tick
			if wheel.Step() {
				for i := 1; i < len(self.wheels); i++ {
					next_wheel := self.wheels[i]
					// degrade timer
					if next_wheel.Step() {
						break
					}
				}
			}
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
	}
}

func (self *TimerEngine) push(timer *Timer) {
	delta := timer.Timestamp - self.start
	if delta < 0 {
		return
	}

	// TODO:无需从零开始
	// TODO:如何Hash还需细致处理
	for i := 0; i < len(self.wheels); i++ {
		wheel := self.wheels[i]
		if delta < wheel.timeMax {
			index := uint64(delta) >> wheel.offset
			wheel.slots[index].PushBack(timer)
			break
		}
	}
}
