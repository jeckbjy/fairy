package timer

import "time"

// Now 当前时间,毫秒单位
func Now() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Start 快速启动定时器,delay时间间隔,period执行次数<=0表示一直执行
func Start(delay int, period int, cb Callback) *Timer {
	if delay <= 0 {
		return nil
	}

	t := New(cb)
	t.delay = delay
	t.period = period
	t.timestamp = Now() + int64(t.delay)
	t.Start()
	return t
}

// StartAt 在某个固定时间执行
func StartAt(timestamp int64, cb Callback) *Timer {
	t := New(cb)
	t.timestamp = timestamp
	t.Start()
	return t
}

// New 创建一个Timer,不启动
func New(cb Callback) *Timer {
	t := &Timer{}
	t.engine = GetEngine()
	t.cb = cb
	return t
}

// Callback 回调函数
type Callback func()

// Timer 定时器,单独线程中执行
type Timer struct {
	prev      *Timer
	next      *Timer
	list      *tlist
	engine    *Engine
	cb        Callback // 回调函数
	timestamp int64    // 到期时间
	delay     int      // 延迟时间
	period    int      // 执行周期
	count     int      // 触发次数
}

// SetEngine 重置engine
func (t *Timer) SetEngine(e *Engine) {
	t.engine = e
}

// Reset 重置时间
func (t *Timer) Reset(ts int64) {
	t.timestamp = ts
}

// Start 启动定时器
func (t *Timer) Start() {
	t.engine.addTimer(t)
}

// Cancel 取消定时器
func (t *Timer) Cancel() {
	t.engine.delTimer(t)
	t.delay = 0
}

func (t *Timer) invoke() {
	t.count++
	t.cb()
	// 继续执行
	if t.delay > 0 && (t.period <= 0 || t.count < t.period) {
		t.timestamp = Now() + int64(t.delay)
		t.Start()
	}
}
