package base

import (
	"container/list"
	"fairy"
)

func NewFilterChain() *BaseFilterChain {
	chain := &BaseFilterChain{}
	chain.New()
	return chain
}

type TravelCallback func(filter fairy.Filter) fairy.FilterAction

type BaseFilterChain struct {
	filters *list.List
}

func (self *BaseFilterChain) New() {
	self.filters = list.New()
}

func (self *BaseFilterChain) AddFirst(filter fairy.Filter) {
	self.filters.PushFront(filter)
}

func (self *BaseFilterChain) AddLast(filter fairy.Filter) {
	self.filters.PushBack(filter)
}

func (self *BaseFilterChain) HandleOpen(conn fairy.Connection) {
	ctx := NewContext(self, conn)
	self.TravelFront(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleOpen(ctx)
	})
}

func (self *BaseFilterChain) HandleClose(conn fairy.Connection) {
	ctx := NewContext(self, conn)
	self.TravelBack(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleClose(ctx)
	})
}

func (self *BaseFilterChain) HandleRead(conn fairy.Connection) {
	// loop read when stop
	ctx := NewContext(self, conn)
	self.TravelFront(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleRead(ctx)
	})
}

func (self *BaseFilterChain) HandleWrite(conn fairy.Connection, msg interface{}) {
	ctx := NewContext(self, conn)
	ctx.SetMessage(msg)
	self.TravelBack(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleWrite(ctx)
	})
}

func (self *BaseFilterChain) HandleError(conn fairy.Connection, err error) {
	ctx := NewContext(self, conn)
	ctx.SetError(err)
	self.TravelFront(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleError(ctx)
	})
}

// process
func (self *BaseFilterChain) TravelFront(cb TravelCallback) {
	for iter := self.filters.Front(); iter != nil; {
		filter := iter.Value.(fairy.Filter)
		action := cb(filter)
		if action == gNextAction {
			iter = iter.Next()
		} else if action == gLastAction {
			iter = self.filters.Back()
		} else {
			break
		}
	}
}

func (self *BaseFilterChain) TravelBack(cb TravelCallback) {
	// 反向遍历
	// for iter := self.filters.Back(); iter != nil {

	// }
}
