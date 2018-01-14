package fairy

import "fairy/util"

var gTimerEngine *TimerEngine

func GetEngine() *TimerEngine {
	if gTimerEngine == nil {
		gTimerEngine = NewEngine()
	}

	return gTimerEngine
}

func NewEngine() *TimerEngine {
	engine := &TimerEngine{}
	engine.executor = GetExecutor()
	return engine
}

type List struct {
	head *Timer
}

func (self *List) Push(timer *Timer) {

}

func (self *List) Pop(timer *Timer) {

}

type Wheel struct {
	slots   []*List
	index   int
	offset  uint
	timeMax int64
}

func (self *Wheel) Tick() bool {
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
	count     int
	wheels    []*Wheel
	executor  *Executor
	pendings  *Timer
}

func (self *TimerEngine) SetExecutor(exec *Executor) {
	self.executor = exec
}

func (self *TimerEngine) AddTimer(timer *Timer) {
	self.count++
	self.push(timer)
}

func (self *TimerEngine) DelTimer(timer *Timer) {
	self.count--
	// owner??
	// del
}

func (self *TimerEngine) Run() {
	go func () {
		// update
	}()
}

func (self *TimerEngine) Update() {
	now := util.Now()
	self.Tick(now)
}

func (self *TimerEngine) Tick(now int64) {
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
			// self.wheels[0].

			// step tick
			if wheel.Tick() {
				for i := 1; i < len(self.wheels); i++ {
					next_wheel := self.wheels[i]
					// degrade timer
					if next_wheel.Tick() {
						break
					}
				}
			}
		}
	}

	// pendings ???
	// self.executor.Dispatch(NewTimerEvent(self))
}

func (self *TimerEngine) Invoke() {
	// invoke all pending timers
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
			wheel.slots[index].Push(timer)
			break
		}
	}
}
