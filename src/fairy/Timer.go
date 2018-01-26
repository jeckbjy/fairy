package fairy

import (
	"fairy/container/inlist"
	"fairy/util"
)

const TIMER_DELAY_MAX = 30 * 24 * 3600 * 1000

type TimerCallback func(*Timer)

type Timer struct {
	inlist.Hook
	engine    *TimerEngine
	cb        TimerCallback
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

func (self *Timer) SetCallback(cb TimerCallback) {
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

func (self *Timer) reset(oldTime int64, newTime int64) {
	if self.Delay <= 0 {
		return
	}
	elapse := self.Timestamp - oldTime
	delay := int64(self.Delay) - elapse
	if delay > 0 {
		self.Delay = int(delay)
		self.Timestamp = newTime + delay
	}
}

func (self *Timer) Restart(timestamp int64) {
	if timestamp < TIMER_DELAY_MAX {
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
			self.engine = GetTimerEngine()
		}
		self.engine.AddTimer(self)
	}
}

func (self *Timer) Stop() {
	if self.IsRunning() {
		self.engine.DelTimer(self)
	}
}

func (self *Timer) IsRunning() bool {
	return self.List() != nil
}

//当timestamp小于TIMER_DELAY_MAX时，代表延迟时间
func NewTimer(timestamp int64, cb TimerCallback, engine *TimerEngine) *Timer {
	delay := 0
	if timestamp < TIMER_DELAY_MAX {
		delay = int(timestamp)
		timestamp = util.Now() + timestamp
	}

	t := &Timer{}
	t.Timestamp = timestamp
	t.Delay = delay
	t.cb = cb
	t.engine = engine
	return t
}

// 快速创建并启动
func StartTimer(timestamp int64, cb TimerCallback) *Timer {
	t := NewTimer(timestamp, cb, GetTimerEngine())
	t.Start()
	return t
}
