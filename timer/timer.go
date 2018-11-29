package timer

import (
	"fmt"
	"time"

	"github.com/jeckbjy/fairy/container/inlist"
)

const (
	ModeTimestamp = 0  // 时间戳执行
	ModeDelay     = 1  // 延迟执行一次
	ModeLoop      = -1 // 循环执行
)

const (
	// Sec 毫秒转秒
	Sec = 1000
)

// TimerCB 定时器回调
type Callback func()

// Start 启动一个定时器
func Start(mode int, ts int64, cb Callback) *Timer {
	t := New(cb)
	t.Start(mode, ts)
	return t
}

// New 创建一个Timer
func New(cb Callback) *Timer {
	t := &Timer{}
	t.cb = cb
	t.engine = GetEngine()
	return t
}

// Timer 定时器
type Timer struct {
	inlist.Hook
	engine    *Engine     // engine
	cb        Callback    // 回调函数
	running   bool        // 是否在运行中
	mode      int         // 类型
	timestamp int64       // 时间戳
	delay     int64       // 不为零，表示是Delay模式
	left      int         // 剩余时间,用于向前调时间时计算
	Count     int         // 记录触发次数
	Tag       int         // 自定义Tag
	Data      interface{} // 自定义数据
}

func (t *Timer) SetEngine(e *Engine) {
	t.Stop()
	t.engine = e
}

func (self *Timer) Start(mode int, ts int64) error {
	if self.isRunning() {
		return fmt.Errorf("timer is running")
	}

	if ts <= 0 {
		return fmt.Errorf("timer bad input")
	}

	self.mode = mode
	self.delay = ts
	self.run()

	return nil
}

// Restart 重新开始,不能是时间戳类型
func (self *Timer) Restart() {
	if self.mode == ModeTimestamp {
		return
	}

	self.Stop()
	self.run()
}

func (self *Timer) Stop() {
	if self.isRunning() {
		self.engine.DelTimer(self)
	}
}

func (self *Timer) run() {
	switch self.mode {
	case ModeTimestamp:
		self.timestamp = self.delay
	case ModeDelay, ModeLoop:
		self.timestamp = time.Now().UnixNano()/int64(time.Millisecond) + self.delay
	}

	self.engine.AddTimer(self)
}

func (self *Timer) isRunning() bool {
	return self.running
}

func (self *Timer) setRunning(flag bool) {
	self.running = flag
}

func (self *Timer) call() {
	if self.cb != nil {
		self.setRunning(false)
		self.Count++
		self.cb()
		// 无线循环
		if self.mode == ModeLoop && self.delay > 0 {
			self.run()
		}
	}
}

// 当向前调时间时,如果是delay模式，会重新计算剩余时间
func (self *Timer) reset(oldTime int64, newTime int64) {
	if self.mode == ModeTimestamp {
		return
	}

	if self.left == 0 {
		self.left = int(self.delay)
	}
	elapse := self.timestamp - oldTime
	delay := int64(self.left) - elapse
	if delay > 0 {
		self.left = int(delay)
		self.timestamp = newTime + delay
	}
}
