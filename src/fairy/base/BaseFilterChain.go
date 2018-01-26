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
	result := self.TravelFront(func(filter fairy.Filter) fairy.FilterAction {
		return filter.HandleError(ctx)
	})

	if result {
		// trigger close when travel all filters
		conn.Close()
	}
}

// TravelFront result:true mean travel all filters
func (self *BaseFilterChain) TravelFront(cb TravelCallback) bool {
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

func (self *BaseFilterChain) TravelBack(cb TravelCallback) bool {
	// 反向遍历
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
