package util

import "sync"

func NewSyncEvent() *SyncEvent {
	ev := &SyncEvent{}
	ev.state = false
	ev.mutex = sync.Mutex{}
	ev.cond = sync.NewCond(&ev.mutex)
	return ev
}

type SyncEvent struct {
	state bool
	cond  *sync.Cond
	mutex sync.Mutex
}

func (self *SyncEvent) Signal() {
	self.mutex.Lock()
	self.state = true
	self.cond.Signal()
	self.mutex.Unlock()
}

func (self *SyncEvent) Broadcast() {
	self.mutex.Lock()
	self.state = true
	self.cond.Broadcast()
	self.mutex.Unlock()
}

func (self *SyncEvent) Wait() {
	self.mutex.Lock()
	for !self.state {
		self.cond.Wait()
	}
	self.mutex.Unlock()
}
