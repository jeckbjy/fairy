package filter

import (
	"fairy"
	"fairy/base"
	"fairy/exec"
	"fairy/log"
	"fairy/rpc"
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
	conn := ctx.GetConn()
	pkt, ok := msg.(fairy.Packet)
	if !ok {
		return ctx.GetNextAction()
	}

	handler := ctx.GetHandler()
	if handler == nil {
		// 通常情况下,在PacketFilter中必然返回了handler,除非自定义了Filter
		if pkt.GetRpcId() > 0 {
			handler = rpc.PopHandler(pkt.GetRpcId())
		}
		if handler == nil {
			handler, _ = self.Dispatcher.GetFinalHandler(pkt.GetId(), pkt.GetName())
		}
	}

	if handler == nil {
		log.Error("cannot find handler:name=%+v,id=%+v", pkt.GetName(), pkt.GetId())
		return ctx.GetStopAction()
	}

	if self.Executor != nil {
		self.DispatchEx(exec.NewPacketEvent(conn, pkt, handler), handler.GetQueueId())
	} else {
		handler.Invoke(conn, pkt)
	}

	return ctx.GetNextAction()
}
