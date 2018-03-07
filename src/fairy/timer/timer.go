package timer

import (
	"fairy/container/inlist"
	"fairy/util"
)

const TIMER_DELAY_MAX = 30 * 24 * 3600 * 1000

type Callback func(*Timer)

type Timer struct {
	inlist.Hook
	engine    *TimerEngine // engine
	cb        Callback     // 回调函数
	Timestamp int64        // 时间戳
	Delay     int          // 不为零，表示是Delay模式
	left      int          // 剩余时间,用于向前调时间时计算
	Count     int          // 记录触发次数
	Tag       int          // 自定义Tag
	Data      interface{}  // 自定义数据
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

func (self *Timer) Call() {
	if self.cb != nil {
		self.Count++
		self.cb(self)
	}
}

// 当向前调时间时,如果是delay模式，会重新计算剩余时间
func (self *Timer) reset(oldTime int64, newTime int64) {
	if self.Delay <= 0 {
		return
	}

	if self.left == 0 {
		self.left = self.Delay
	}
	elapse := self.Timestamp - oldTime
	delay := int64(self.left) - elapse
	if delay > 0 {
		self.left = int(delay)
		self.Timestamp = newTime + delay
	}
}

func (self *Timer) Restart() {
	if self.Delay != 0 {
		self.Start(int64(self.Delay))
	}
}

func (self *Timer) Start(ts int64) {
	if !self.IsRunning() {
		if ts < TIMER_DELAY_MAX {
			self.Delay = int(ts)
			self.left = self.Delay
			ts = util.Now() + ts
		}

		self.Timestamp = ts
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

func Sec(t int64) int64 {
	return t * 1000
}

//当timestamp小于TIMER_DELAY_MAX时，代表延迟时间
func New(engine *TimerEngine, cb Callback) *Timer {
	if engine == nil {
		engine = GetTimerEngine()
	}

	t := &Timer{}
	t.cb = cb
	t.engine = engine
	return t
}

// 快速创建并启动
func Start(ts int64, cb Callback) *Timer {
	t := New(nil, cb)
	t.Start(ts)
	return t
}

// 异步，不需要post主线程中执行
func StartAsync(ts int64, cb Callback) *Timer {
	t := New(GetAsyncTimerEngine(), cb)
	t.Start(ts)
	return t
}
