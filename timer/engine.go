package timer

import (
	"sync"
	"time"
)

var gEngine = NewEngine()

// GetEngine get global timer engine
func GetEngine() *Engine {
	return gEngine
}

// NewEngine create timer engine
func NewEngine() *Engine {
	e := &Engine{}
	e.init()
	return e
}

// Engine Hashed and Hierarchical Timing Wheels
// thread safe
type Engine struct {
	stopped   bool       // 是否在执行
	start     int64      // 起始时间
	timestamp int64      // 当前时间戳
	count     int        // timer个数
	wheels    []*twheel  // 时间桶
	mutex     sync.Mutex // 锁
	timeMax   uint64     // 当前最大
	interval  int        // 间隔
}

func (e *Engine) init() {
	e.stopped = true
	e.timestamp = time.Now().UnixNano() / int64(time.Millisecond)
	e.start = e.timestamp
	e.count = 0
	e.interval = cfgTimeInterval

	// create wheel
	for i := 0; i < cfgWheelNum; i++ {
		e.newWheel()
	}

	go e.Start()
}

func (e *Engine) newWheel() {
	wheel := &twheel{}
	wheel.init(len(e.wheels))
	e.timeMax = wheel.timeMax
	e.wheels = append(e.wheels, wheel)
}

func (e *Engine) loop() {
	for !e.stopped {
		e.Update()
		time.Sleep(time.Duration(e.interval) * time.Millisecond)
	}
}

func (e *Engine) Start() {
	if e.stopped {
		e.stopped = false
		go e.loop()
	}
}

func (e *Engine) Stop() {
	e.stopped = true
}

// Update 以当前时间更新
func (e *Engine) Update() {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	e.Tick(timestamp)
}

// Tick 检测定时器
func (e *Engine) Tick(newTime int64) {
	oldTime := e.timestamp
	e.timestamp = newTime

	pendings := tlist{}
	e.mutex.Lock()

	if newTime < oldTime {
		// 向前调时间导致时间变小
		// 收集所有Timer并重新计算
		for _, wheel := range e.wheels {
			wheel.index = 0
			for _, slot := range wheel.slots {
				e.calcTimers(&pendings, slot)
			}
		}
	} else {
		// 正常tick
		ticks := newTime - oldTime
		for i := int64(0); i < ticks; i++ {
			// wheels[0][0] 已经到时间
			e.calcTimers(&pendings, e.wheels[0].current())
			// cascade 向前走一个刻度
			for i := 0; i < len(e.wheels); i++ {
				if !e.wheels[i].step() {
					// 还没有走一圈
					break
				}

				if i+1 == len(e.wheels) {
					// 溢出,新建一个
					e.newWheel()
					break
				}

				e.calcTimers(&pendings, e.wheels[i+1].current())
			}
		}
	}

	e.mutex.Unlock()

	// 执行timer
	for node := pendings.head; node != nil; {
		curr := node
		// TODO:这里可能会有线程安全问题,会修改timer的list值
		node = pendings.remove(node)
		curr.invoke()
	}
}

func (e *Engine) calcTimers(pendings *tlist, slot *tlist) {
	for node := slot.head; node != nil; {
		curr := node
		node = slot.remove(node)
		if curr.timestamp < e.timestamp {
			pendings.push(curr)
		} else {
			e.push(curr)
		}
	}
}

func (e *Engine) push(t *Timer) {
	delta := uint64(t.timestamp - e.timestamp)

	for i := 0; i < len(e.wheels); i++ {
		wheel := e.wheels[i]
		if delta < wheel.timeMax {
			wheel.push(t, delta)
			break
		}
	}
}

// addTimer 添加定时器
func (e *Engine) addTimer(t *Timer) {
	e.mutex.Lock()
	if t.list == nil && t.timestamp > e.timestamp {
		e.count++
		e.push(t)
	}
	e.mutex.Unlock()
}

// delTimer 删除定时器
func (e *Engine) delTimer(t *Timer) {
	e.mutex.Lock()
	if t.list != nil {
		t.list.remove(t)
		e.count--
	}
	e.mutex.Unlock()
}
