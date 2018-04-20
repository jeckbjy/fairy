package filter

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
	"github.com/jeckbjy/fairy/log"
	"github.com/jeckbjy/fairy/rpc"
)

func NewExecutor() *zExecutorFilter {
	return NewExecutorEx(fairy.GetExecutor(), fairy.GetDispatcher())
}

func NewExecutorEx(e *fairy.Executor, d *fairy.Dispatcher) *zExecutorFilter {
	filter := &zExecutorFilter{}
	filter.Executor = e
	filter.Dispatcher = d
	return filter
}

// 默认的执行线程
type zExecutorFilter struct {
	base.BaseFilter
	*fairy.Executor
	*fairy.Dispatcher
}

func (self *zExecutorFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
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
		self.DispatchEx(fairy.NewPacketEvent(conn, pkt, handler), handler.GetQueueId())
	} else {
		handler.Invoke(conn, pkt)
	}

	return ctx.GetNextAction()
}
