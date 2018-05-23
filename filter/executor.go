package filter

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

func NewExecutor() *ExecutorFilter {
	filter := &ExecutorFilter{Executor: fairy.GetExecutor()}
	return filter
}

// 默认的执行线程
type ExecutorFilter struct {
	base.BaseFilter
	*fairy.Executor
}

func (self *ExecutorFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	data := ctx.GetMessage()
	handlerCtx, ok := data.(*fairy.HandlerCtx)
	if !ok {
		return ctx.GetNextAction()
	}

	if self.Executor != nil {
		self.DispatchEx(handlerCtx, handlerCtx.QueueID())
	} else {
		handlerCtx.Process()
	}

	return ctx.GetNextAction()
}
