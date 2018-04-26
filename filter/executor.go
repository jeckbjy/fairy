package filter

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

func NewExecutor() *zExecutorFilter {
	filter := &zExecutorFilter{Executor: fairy.GetExecutor()}
	return filter
}

// 默认的执行线程
type zExecutorFilter struct {
	base.BaseFilter
	*fairy.Executor
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
		// log
		return ctx.GetStopAction()
	}

	if self.Executor != nil {
		self.DispatchEx(fairy.NewPacketEvent(conn, pkt, handler), handler.QueueID())
	} else {
		ctx := fairy.HandlerCtx{Conn: conn, Packet: pkt}
		handler.Invoke(&ctx)
	}

	return ctx.GetNextAction()
}
