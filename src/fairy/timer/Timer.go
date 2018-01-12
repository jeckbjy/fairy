package timer

import (
	"fairy/util"
)

/*
example 1:
func OnTimeout(t *timer.Timer) {
	fmt.Println("timeout")
}
1:
t := timer.New(10, true, OnTimeout)
t.Start()
2:
timer.Start(util.Now() + 10, OnTimeout)
3:
timer.StartDelay(10, OnTimeout)

example 2:
func OnTimeout(t *timer.Timer) {
	if t.Count < 10 {
		t.Restart(10, true)
	}
}

func Start() {
	//timer.Start(10, OnTimeout)
	t := timer.New(10, true, OnTimeout)
	t.Tag = 1
	t.Data = 1
	t.Start()
}
*/
type Callback func(*Timer)

type Timer struct {
	owner     *List
	prev      *Timer
	next      *Timer
	engine    *TimerEngine
	cb        Callback
	Timestamp int64       // 时间戳
	Delay     int         // 不为零，表示是Delay模式
	Count     int         // 触发次数
	Tag       int         // 自定义Tag
	Data      interface{} // 自定义数据
}

func (self *Timer) SetEngine(e *TimerEngine) {
	self.Stop()
	self.engine = e
}

func (self *Timer) SetCallback(cb Callback) {
	self.cb = cb
}

func (self *Timer) SetTag(tag int) {
	self.Tag = tag
}

func (self *Timer) SetData(data interface{}) {
	self.Data = data
}

func (self *Timer) Invoke() {
	if self.cb != nil {
		self.Count++
		self.cb(self)
	}
}

func (self *Timer) IsDelayMode() bool {
	return self.Delay != 0
}

func (self *Timer) Restart(timestamp int64, delayMode bool) {
	if delayMode {
		self.Delay = int(timestamp)
		timestamp = util.Now() + timestamp
	}

	self.Timestamp = timestamp
	self.Stop()
	self.Start()
}

func (self *Timer) Start() {
	if !self.IsRunning() {
		if self.engine == nil {
			self.engine = GetEngine()
		}
		// check engine running??
		self.engine.AddTimer(self)
	}
}

func (self *Timer) Stop() {
	if self.IsRunning() {
		self.engine.DelTimer(self)
	}
}

func (self *Timer) IsRunning() bool {
	return self.prev != nil
}

/*
Create Timer Func
*/
func New(timestamp int64, delayMode bool, cb Callback) *Timer {
	delay := 0
	if delayMode {
		delay = int(timestamp)
		timestamp = util.Now() + timestamp
	}

	t := &Timer{}
	t.Timestamp = timestamp
	t.Delay = delay
	t.cb = cb
	return t
}

func Start(timestamp int64, cb Callback) *Timer {
	t := New(timestamp, false, cb)
	t.Start()
	return t
}

func StartDelay(delay int, cb Callback) *Timer {
	t := New(int64(delay), true, cb)
	t.Start()
	return t
}
