package fairy

import (
	"container/list"
	"sync"
)

func NewEventQueue() *EventQueue {
	eq := &EventQueue{}
	return eq
}

type EventQueue struct {
	events  *list.List
	mutex   sync.Mutex
	stopped bool
}

func (self *EventQueue) Push(ev Event) {
	self.mutex.Lock()
	self.events.PushBack(ev)
	self.mutex.Unlock()
	// 通知执行
}

func (self *EventQueue) loop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for self.stopped {

	}
}

func (self *EventQueue) Start(wg *sync.WaitGroup) {
	go self.loop(wg)
}

func (self *EventQueue) Stop() {
	self.stopped = true
}
