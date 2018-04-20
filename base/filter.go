package base

import "github.com/jeckbjy/fairy"

// BaseFilter 实现所有接口
type BaseFilter struct {
}

func (filter *BaseFilter) Stateless() bool {
	return true
}

func (self *BaseFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	return ctx.GetNextAction()
}

func (self *BaseFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	return ctx.GetNextAction()
}

func (self *BaseFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	return ctx.GetNextAction()
}

func (self *BaseFilter) HandleClose(ctx fairy.FilterContext) fairy.FilterAction {
	return ctx.GetNextAction()
}

func (self *BaseFilter) HandleError(ctx fairy.FilterContext) fairy.FilterAction {
	return ctx.GetNextAction()
}
