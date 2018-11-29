package timer

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/jeckbjy/fairy/exit"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/container/inlist"
)

var gEngine *Engine

func init() {
	gEngine = NewEngine(fairy.GetExecutor())
	gEngine.Start()
}

// GetEngine 同步主线程定时器
func GetEngine() *Engine {
	return gEngine
}

// NewEngine 创建一个Engine
func NewEngine(executor *fairy.Executor) *Engine {
	engine := &Engine{}
	engine.create(executor)
	return engine
}

// Engine use Hashed and Hierarchical Timing Wheels
type Engine struct {
	start     int64           // 起始时间
	timestamp int64           // 当前时间戳
	interval  int64           // 最小触发时间间隔
	count     int32           // timer个数
	wheels    []*Wheel        // 桶
	executor  *fairy.Executor //
	pendings  *inlist.List
	mutex     *sync.Mutex
	stopped   bool
	timeMax   uint64
}

func (self *Engine) create(exec *fairy.Executor) {
	self.SetInterval(TIME_INTERVAL)
	self.executor = exec
	self.timestamp = time.Now().UnixNano() / int64(time.Millisecond)
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

func (self *Engine) newWheel() {
	wheel := &Wheel{}
	wheel.Create(len(self.wheels))
	self.timeMax = wheel.timeMax
	self.wheels = append(self.wheels, wheel)
}

func (self *Engine) SetExecutor(executor *fairy.Executor) {
	self.executor = executor
}

func (self *Engine) SetInterval(interval int) {
	self.interval = int64(interval)
}

func (self *Engine) AddTimer(timer *Timer) {
	if timer.timestamp > self.timestamp {
		self.mutex.Lock()
		atomic.AddInt32(&self.count, 1)
		timer.setRunning(true)
		self.push(timer)
		self.mutex.Unlock()
	}
}

func (self *Engine) DelTimer(timer *Timer) {
	self.mutex.Lock()
	if timer.isRunning() {
		atomic.AddInt32(&self.count, -1)
		timer.setRunning(false)
		timer.List().Remove(timer)
	}
	self.mutex.Unlock()
}

func (self *Engine) Run() {
	for !self.stopped {
		self.Update()
		time.Sleep(time.Duration(self.interval) * time.Millisecond)
	}
}

func (self *Engine) Start() {
	if self.stopped {
		exit.Add(self.Stop)
		self.stopped = false
		go self.Run()
	}
}

func (self *Engine) Stop() {
	if !self.stopped {
		self.stopped = true
	}
}

func (self *Engine) Update() {
	newTime := time.Now().UnixNano() / int64(time.Millisecond)
	self.UpdateBy(newTime)
}

func (self *Engine) UpdateBy(newTime int64) {
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
		self.executor.Dispatch(0, func() {
			self.Invoke()
		})
	}

	self.mutex.Unlock()

	// async invoke if no exector
	if pendings.Len() > 0 {
		self.triggerAllTimer(&pendings)
	}
}

func (self *Engine) Invoke() {
	// invoke all pending timers
	pendings := inlist.List{}
	self.mutex.Lock()
	pendings.MoveBackList(self.pendings)
	self.pendings.Init()
	self.mutex.Unlock()

	// invoke
	self.triggerAllTimer(&pendings)
}

func (self *Engine) triggerAllTimer(pendings *inlist.List) {
	for iter := pendings.Front(); iter != nil; {
		timer := iter.(*Timer)
		iter = inlist.Next(iter)
		pendings.Remove(timer)
		atomic.AddInt32(&self.count, -1)
		timer.call()
	}
}

func (self *Engine) rebuild(pendings *inlist.List, oldTime int64, newTime int64) {
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
		if timer.timestamp >= newTime {
			pendings.Remove(timer)
			self.push(timer)
		}
	}
}

func (self *Engine) tick(pendings *inlist.List, oldTime int64, newTime int64) {
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

func (self *Engine) cascade(pendings *inlist.List) {
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

			if timer.timestamp <= self.timestamp {
				pendings.PushBack(timer)
			} else {
				self.push(timer)
			}
		}
	}
}

func (self *Engine) push(timer *Timer) {
	if timer.timestamp <= self.timestamp {
		// panic("bad timer")
		timer.setRunning(false)
		atomic.AddInt32(&self.count, -1)
		return
	}

	delta := uint64(timer.timestamp - self.timestamp)
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
