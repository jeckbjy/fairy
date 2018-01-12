package filter

import (
	"fairy"
	"fairy/base"
)

func NewExecutorFilter() *ExecutorFilter {
	return NewExecutorFilterEx(fairy.GetExecutor(), fairy.GetDispatcher())
}

func NewExecutorFilterEx(e *fairy.Executor, d *fairy.Dispatcher) *ExecutorFilter {
	filter := &ExecutorFilter{}
	filter.Executor = e
	filter.Dispatcher = d
	return filter
}

// 默认的执行线程
type ExecutorFilter struct {
	base.BaseFilter
	*fairy.Executor
	*fairy.Dispatcher
}

func (self *ExecutorFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	msg := ctx.GetMessage()
	conn := ctx.GetConnection()
	if packet, ok := msg.(fairy.Packet); ok {
		if invoker := self.GetHandler(packet.GetId(), packet.GetName()); invoker != nil {
			if self.Executor != nil {
				self.DispatchEx(fairy.NewPacketEvent(conn, packet, invoker), invoker.GetQueueId())
			} else {
				invoker.Invoke(conn, packet)
			}
		}
	}
	return ctx.GetNextAction()
}
