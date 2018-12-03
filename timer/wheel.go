package timer

const (
	cfgTimeInterval = 1               // 默认时间间隔
	cfgWheelNum     = 3               // 初始wheel个数(2^30)
	cfgSlotPow      = 10              //
	cfgSlotMax      = 1 << cfgSlotPow //
)

// tlist 双向非循环链表
type tlist struct {
	head *Timer
	tail *Timer
}

func (l *tlist) push(t *Timer) {
	if l.head == nil {
		l.head = t
		l.tail = t
	} else {
		t.prev = l.tail
		l.tail = t
	}

	t.list = l
}

func (l *tlist) remove(t *Timer) *Timer {
	next := t.next
	if t.prev != nil {
		t.prev.next = t.next
	}

	if t.next != nil {
		t.next.prev = t.prev
	}

	if t == l.head {
		l.head = t.next
	}

	if t == l.tail {
		l.tail = t.prev
	}

	t.list = nil
	return next
}

// twheel timer wheel
type twheel struct {
	slots   []*tlist // timer slot
	index   int      // 当前slot索引
	shift   uint     // shift offset
	timeMax uint64   // timestamp max
}

func (w *twheel) init(index int) {
	w.index = 0
	w.shift = uint(index * cfgSlotPow)
	w.timeMax = uint64(1) << uint((index+1)*cfgSlotPow)
	for i := 0; i < cfgSlotMax; i++ {
		w.slots = append(w.slots, &tlist{})
	}
}

func (w *twheel) current() *tlist {
	return w.slots[w.index]
}

// step 前进一步,返回是否到了一圈
func (w *twheel) step() bool {
	w.index++
	if w.index >= len(w.slots) {
		w.index = 0
		return true
	}

	return false
}

func (w *twheel) push(t *Timer, delta uint64) {
	off := (delta >> w.shift) - 1
	index := (int(off) + w.index) % len(w.slots)
	w.slots[index].push(t)
}
