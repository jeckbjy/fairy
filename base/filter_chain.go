package base

import (
	"container/list"

	"github.com/jeckbjy/fairy"
)

func NewFilterChain() *FilterChain {
	chain := &FilterChain{}
	chain.New()
	return chain
}

// ScanCB scan filter callback
type ScanCB func(filter fairy.Filter) fairy.FilterAction

type FilterChain struct {
	filters *list.List
}

func (self *FilterChain) New() {
	self.filters = list.New()
}

func (self *FilterChain) Len() int {
	return self.filters.Len()
}

func (self *FilterChain) AddFirst(filter fairy.Filter) {
	self.filters.PushFront(filter)
}

func (self *FilterChain) AddLast(filter fairy.Filter) {
	self.filters.PushBack(filter)
}

func (self *FilterChain) HandleOpen(conn fairy.Conn) {
	ctx := NewContext(self, conn)
	self.ScanFront(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleOpen(ctx)
	})
}

func (self *FilterChain) HandleClose(conn fairy.Conn) {
	ctx := NewContext(self, conn)
	self.ScanBack(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleClose(ctx)
	})
}

func (self *FilterChain) HandleRead(conn fairy.Conn) {
	// loop read when stop
	ctx := NewContext(self, conn)
	self.ScanFront(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleRead(ctx)
	})
}

func (self *FilterChain) HandleWrite(conn fairy.Conn, msg interface{}) {
	ctx := NewContext(self, conn)
	ctx.SetMessage(msg)
	self.ScanBack(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleWrite(ctx)
	})
}

func (self *FilterChain) HandleError(conn fairy.Conn, err error) {
	ctx := NewContext(self, conn)
	ctx.SetError(err)
	result := self.ScanFront(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleError(ctx)
	})

	if result {
		// trigger close when travel all filters
		conn.Close()
	}
}

// ScanFront result:true mean travel all filters
func (self *FilterChain) ScanFront(cb ScanCB) bool {
	for iter := self.filters.Front(); iter != nil; {
		filter := iter.Value.(fairy.Filter)
		action := cb(filter)
		switch action {
		case gNextAction:
			iter = iter.Next()
		case gLastAction:
			iter = self.filters.Back()
		case gFirstAction:
			iter = self.filters.Front()
		case gStopAction:
			return false
		default:
			return false
		}
	}

	return true
}

// ScanBack 反向遍历
func (self *FilterChain) ScanBack(cb ScanCB) bool {
	for iter := self.filters.Back(); iter != nil; {
		filter := iter.Value.(fairy.Filter)
		action := cb(filter)
		switch action {
		case gNextAction:
			iter = iter.Prev()
		case gLastAction:
			iter = self.filters.Front()
		case gFirstAction:
			iter = self.filters.Back()
		case gStopAction:
			return false
		default:
			return false
		}
	}

	return true
}
