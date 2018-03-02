package base

import (
	"container/list"
	"fairy"
	"io"
	"sync"
)

type ConnWriter struct {
	buffers *list.List
	cond    *sync.Cond
	mutex   *sync.Mutex
	future  *BaseFuture
	stopped bool
}

func (self *ConnWriter) NewWriter() {
	self.mutex = &sync.Mutex{}
	self.stopped = true
}

func (self *ConnWriter) lazyInit() {
	if self.buffers == nil {
		self.buffers = list.New()
		self.cond = sync.NewCond(self.mutex)
		self.future = NewFuture()
	}
}

func (self *ConnWriter) IsStopped() bool {
	return self.stopped
}

func (self *ConnWriter) StopWrite() {
	if !self.stopped {
		self.stopped = true
		self.cond.Signal()
	}
}

func (self *ConnWriter) Flush() {
	if self.future != nil {
		self.future.Wait(-1)
	}
}

func (self *ConnWriter) PushBuffer(buffer *fairy.Buffer, cb fairy.Callback) {
	self.mutex.Lock()
	self.lazyInit()
	self.buffers.PushBack(buffer)
	self.future.Reset()
	if self.stopped {
		self.stopped = false
		go cb()
	}
	self.cond.Signal()
	self.mutex.Unlock()
}

func (self *ConnWriter) WaitBuffers(buffers *list.List) {
	self.mutex.Lock()
	for !self.stopped && self.buffers.Len() == 0 {
		self.future.DoneSucceed()
		self.cond.Wait()
	}

	*buffers = *self.buffers
	self.buffers.Init()
	self.mutex.Unlock()
}

func (self *ConnWriter) WriteBuffers(writer io.Writer, l *list.List) error {
	for iterl := l.Front(); iterl != nil; iterl = iterl.Next() {
		if iterl.Value == nil {
			return nil
		}
		buffer := iterl.Value.(*fairy.Buffer)
		for iterb := buffer.Front(); iterb != nil; iterb = iterb.Next() {
			data := iterb.Value.([]byte)
			_, err := writer.Write(data)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}
