package timer

import "github.com/jeckbjy/fairy/container/inlist"

const (
	TIME_INTERVAL = 1             // 默认时间间隔
	WHEEL_NUM     = 3             // 初始wheel个数(2^30)，可以扩展
	SLOT_POW      = 10            // 2^SLOT_POW,默认10
	SLOT_MAX      = 1 << SLOT_POW // 个数
)

type Wheel struct {
	slots   []*inlist.List // 桶
	index   int            // 当前slot循环索引
	timeOff uint           // shift offset
	timeMax uint64         // 区间最大值
}

func (self *Wheel) Create(index int) {
	self.index = 0
	self.timeOff = uint(index * SLOT_POW)
	self.timeMax = uint64(1) << uint((index+1)*SLOT_POW)
	for i := 0; i < SLOT_MAX; i++ {
		self.slots = append(self.slots, inlist.New())
	}
}

func (self *Wheel) Current() *inlist.List {
	return self.slots[self.index]
}

func (self *Wheel) Step() bool {
	self.index++
	if self.index >= len(self.slots) {
		self.index = 0
		return true
	}

	return false
}

func (self *Wheel) Push(timer *Timer, delta uint64) {
	off := (delta >> self.timeOff) - 1
	index := (int(off) + self.index) % len(self.slots)
	self.slots[index].PushBack(timer)
}
