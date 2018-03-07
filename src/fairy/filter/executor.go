package filter

import (
	"fairy"
	"fairy/base"
	"fairy/exec"
	"fairy/log"
)

func NewExecutor() *ExecutorFilter {
	return NewExecutorEx(exec.GetExecutor(), fairy.GetDispatcher())
}

func NewExecutorEx(e *exec.Executor, d *fairy.Dispatcher) *ExecutorFilter {
	filter := &ExecutorFilter{}
	filter.Executor = e
	filter.Dispatcher = d
	return filter
}

// 默认的执行线程
type ExecutorFilter struct {
	base.BaseFilter
	*exec.Executor
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
		// ctx.ThrowError(fmt)
		log.Error("cannot find handler:name=%+v,id=%+v", packet.GetName(), packet.GetId())
		return ctx.GetStopAction()
	}

	if self.Executor != nil {
		self.DispatchEx(exec.NewPacketEvent(conn, packet, handler), handler.GetQueueId())
	} else {
		handler.Invoke(conn, packet)
	}

	return ctx.GetNextAction()
}
