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
	packet, ok := msg.(fairy.Packet)
	if !ok {
		return ctx.GetNextAction()
	}

	// check rpc
	if packet.GetSerialId() > 0 {

	}

	// check handler
	handler := self.GetHandler(packet.GetId(), packet.GetName())
	if handler == nil {
		handler = self.GetUncaughtHandler()
	}

	if handler == nil {
		// TODO:throw error
		fairy.Error("cannot find handler:name=%+v,id=%+v", packet.GetName(), packet.GetId())
		return ctx.GetStopAction()
	}

	if self.Executor != nil {
		self.DispatchEx(fairy.NewPacketEvent(conn, packet, handler), handler.GetQueueId())
	} else {
		defer fairy.Catch()
		handler.Invoke(conn, packet)
	}

	return ctx.GetNextAction()
}
