package filters

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

func NewExecutor() *ExecutorFilter {
	filter := &ExecutorFilter{}
	return filter
}

// ExecutorFilter 默认的消息执行
type ExecutorFilter struct {
	base.Filter
}

func (ef *ExecutorFilter) Name() string {
	return "ExecutorFilter"
}

// HandleRead 转发到调度器中执行
func (ef *ExecutorFilter) HandleRead(ctx fairy.IFilterCtx) {
	data := ctx.GetData()

	if pkt, ok := data.(fairy.IPacket); ok {
		handler := fairy.GetDispatcher().GetHandler(pkt.GetId(), pkt.GetName())
		if handler != nil {
			fairy.GetExecutor().Dispatch(handler.QueueId, func() {
				fairy.InvokeHandler(ctx.GetConn(), pkt, handler)
			})
		}
	}

	ctx.Next()
}
