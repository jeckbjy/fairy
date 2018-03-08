package timer

import (
	"fairy/container/inlist"
	"fairy/exec"
	"fairy/util"
	"sync"
	"sync/atomic"
	"time"
)

var gAsyncTimerEngine *TimerEngine
var gTimerEngine *TimerEngine

// 异步定时器线程，不会在主线程中执行
func GetAsyncTimerEngine() *TimerEngine {
	util.Once(gAsyncTimerEngine, func() {
		gAsyncTimerEngine = NewEngine(nil)
		gAsyncTimerEngine.Start()
	})

	return gAsyncTimerEngine
}

// 同步主线程定时器
func GetTimerEngine() *TimerEngine {
	util.OnceEx(gTimerEngine, func() {
		gTimerEngine = NewEngine(exec.GetExecutor())
		gTimerEngine.Start()
	})

	return gTimerEngine
}

func NewEngine(executor *exec.Executor) *TimerEngine {
	engine := &TimerEngine{}
	engine.Create(executor)
	return engine
}

type TimerEvent struct {
	engine *TimerEngine
}

func (self *TimerEvent) Process() {
	self.engine.Invoke()
}

type TimerEngine struct {
	start     int64
	timestamp int64
	interval  int64
	count     int32
	wheels    []*Wheel
	executor  *exec.Executor
	pendings  *inlist.List
	mutex     *sync.Mutex
	stopped   bool
	timeMax   uint64
}

func (self *TimerEngine) Create(exec *exec.Executor) {
	self.SetInterval(TIME_INTERVAL)
	self.executor = exec
	self.timestamp = util.Now()
	self.start = self.timestamp
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

func (self *TimerEngine) SetExecutor(executor *exec.Executor) {
	self.executor = executor
}

func (self *TimerEngine) SetInterval(interval int) {
	self.interval = int64(interval)
}

func (self *TimerEngine) AddTimer(timer *Timer) {
	if timer.Timestamp > self.timestamp {
		self.mutex.Lock()
		atomic.AddInt32(&self.count, 1)
		timer.setRunning(true)
		self.push(timer)
		self.mutex.Unlock()
	}
}

func (self *TimerEngine) DelTimer(timer *Timer) {
	self.mutex.Lock()
	if timer.isRunning() {
		atomic.AddInt32(&self.count, -1)
		timer.setRunning(false)
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
		util.RegisterExit(self)
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
	newTime := util.Now()
	self.UpdateBy(newTime)
}

func (self *TimerEngine) UpdateBy(newTime int64) {
	oldTime := self.timestamp
	self.timestamp = newTime
	self.mutex.Lock()
	pendings := inlist.List{}

	if newTime <= oldTime {
		self.rebuild(&pendings, oldTime, newTime)
	} else {
		self.tick(&pendings, oldTime, newTime)
	}

	if self.executor != nil && pendings.Len() > 0 {
		self.pendings.MoveBackList(&pendings)
		self.executor.Dispatch(&TimerEvent{engine: self})
	}

	self.mutex.Unlock()

	// async call
	for iter := pendings.Front(); iter != nil; {
		timer := iter.(*Timer)
		iter = inlist.Next(iter)
		pendings.Remove(timer)
		timer.setRunning(false)
		atomic.AddInt32(&self.count, -1)
		timer.Call()
	}
}

func (self *TimerEngine) rebuild(pendings *inlist.List, oldTime int64, newTime int64) {
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
}

func (self *TimerEngine) tick(pendings *inlist.List, oldTime int64, newTime int64) {
	// tick all timer
	ticks := newTime - oldTime
	for i := int64(0); i < ticks; i++ {
		wheel := self.wheels[0]
		// process timer list
		pendings.MoveBackList(wheel.Current())
		// step tick
		self.cascade(pendings)
	}
}

func (self *TimerEngine) cascade(pendings *inlist.List) {
	for i := 0; i < len(self.wheels); i++ {
		if !self.wheels[i].Step() {
			break
		}

		if i+1 == len(self.wheels) {
			self.newWheel()
			break
		}

		// rehash next wheel
		slots := self.wheels[i+1].Current()
		for iter := slots.Front(); iter != nil; {
			timer := iter.(*Timer)
			iter = inlist.Next(iter)
			slots.Remove(timer)

			if timer.Timestamp <= self.timestamp {
				pendings.PushBack(timer)
			} else {
				self.push(timer)
			}
		}
	}
}

func (self *TimerEngine) Invoke() {
	// invoke all pending timers
	pendings := inlist.List{}
	self.mutex.Lock()
	pendings.MoveBackList(self.pendings)
	self.pendings.Init()
	self.mutex.Unlock()

	// invoke
	for iter := pendings.Front(); iter != nil; {
		timer := iter.(*Timer)
		iter = inlist.Next(iter)
		pendings.Remove(timer)
		atomic.AddInt32(&self.count, -1)
		timer.setRunning(false)
		timer.Call()
	}
}

func (self *TimerEngine) push(timer *Timer) {
	if timer.Timestamp <= self.timestamp {
		// panic("bad timer")
		timer.setRunning(false)
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
