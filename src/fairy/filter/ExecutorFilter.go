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
		handler := self.GetHandler(packet.GetId(), packet.GetName())
		if handler == nil {
			handler = self.GetUncaughtHandler()
		}

		if handler != nil {
			if self.Executor != nil {
				self.DispatchEx(fairy.NewPacketEvent(conn, packet, handler), handler.GetQueueId())
			} else {
				handler.Invoke(conn, packet)
			}
		} else {
			fairy.Error("cannot find handler:name=%+v,id=%+v", packet.GetName(), packet.GetId())
		}
	}
	return ctx.GetNextAction()
}
