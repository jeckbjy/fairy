package fairy

import (
	"container/list"
	"fairy/util"
	"sync"
)

func NewEventQueue() *EventQueue {
	eq := &EventQueue{}
	eq.stopped = true
	eq.mutex = &sync.Mutex{}
	eq.cond = sync.NewCond(eq.mutex)
	eq.events = list.New()
	return eq
}

type EventQueue struct {
	events  *list.List
	mutex   *sync.Mutex
	cond    *sync.Cond
	stopped bool
}

func (self *EventQueue) Push(ev Event) {
	self.mutex.Lock()
	self.events.PushBack(ev)
	// notify process
	self.cond.Signal()
	self.mutex.Unlock()
}

func (self *EventQueue) loop(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		events := list.List{}
		self.mutex.Lock()
		for !self.stopped && self.events.Len() == 0 {
			self.cond.Wait()
		}
		util.SwapList(&events, self.events)
		self.mutex.Unlock()
		// process all events
		for iter := events.Front(); iter != nil; iter = iter.Next() {
			ev := iter.Value.(Event)
			ev.Process()
		}
	}
}

func (self *EventQueue) Start(wg *sync.WaitGroup) {
	if self.stopped {
		self.stopped = false
		go self.loop(wg)
	}
}

func (self *EventQueue) Stop() {
	self.mutex.Lock()
	self.stopped = true
	self.cond.Signal()
	self.mutex.Unlock()
}
