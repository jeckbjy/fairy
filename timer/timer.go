package timer

import (
	"github.com/jeckbjy/fairy/container/inlist"
	"github.com/jeckbjy/fairy/util"
)

const (
	// Sec 毫秒转秒
	Sec = 1000
)

// New 当timestamp小于TIMER_DELAY_MAX时，代表延迟时间
func New(engine *TimerEngine, cb TimerCB) *Timer {
	if engine == nil {
		engine = GetEngine()
	}

	t := &Timer{}
	t.cb = cb
	t.engine = engine
	return t
}

// Start 快速创建并启动
func Start(ts int64, cb TimerCB) *Timer {
	t := New(nil, cb)
	t.Start(ts)
	return t
}

// StartAsync 异步，不需要post主线程中执行
func StartAsync(ts int64, cb TimerCB) *Timer {
	t := New(GetAsyncEngine(), cb)
	t.Start(ts)
	return t
}

// 大于zTimerDelayMax表示时间戳
const zTimerDelayMax = 30 * 24 * 3600 * 1000

// TimerCB 定时器回调函数
type TimerCB func(*Timer)

// Timer 定时器
type Timer struct {
	inlist.Hook
	engine    *TimerEngine // engine
	cb        TimerCB      // 回调函数
	Timestamp int64        // 时间戳
	Delay     int          // 不为零，表示是Delay模式
	left      int          // 剩余时间,用于向前调时间时计算
	Count     int          // 记录触发次数
	Tag       int          // 自定义Tag
	Data      interface{}  // 自定义数据
	running   bool         // 是否在运行中
}

func (self *Timer) SetEngine(e *TimerEngine) {
	self.Stop()
	self.engine = e
}

func (self *Timer) SetCallback(cb TimerCB) {
	self.cb = cb
}

func (self *Timer) SetTag(tag int) {
	self.Tag = tag
}

func (self *Timer) SetData(data interface{}) {
	self.Data = data
}

func (self *Timer) Restart() {
	if self.Delay != 0 {
		self.Start(int64(self.Delay))
	}
}

func (self *Timer) Start(ts int64) {
	if !self.isRunning() {
		if ts < zTimerDelayMax {
			self.Delay = int(ts)
			self.left = self.Delay
			ts = util.Now() + ts
		}

		self.Timestamp = ts
		self.engine.AddTimer(self)
	}
}

func (self *Timer) Stop() {
	if self.isRunning() {
		self.engine.DelTimer(self)
	}
}

func (self *Timer) isRunning() bool {
	return self.running
}

func (self *Timer) setRunning(flag bool) {
	self.running = flag
}

func (self *Timer) call() {
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
