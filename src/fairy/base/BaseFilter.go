package base

import (
	. "fairy"
)

// 实现所有接口
type BaseFilter struct {
}

func (self *BaseFilter) HandleRead(ctx FilterContext) FilterAction {
	return ctx.GetNextAction()
}

func (self *BaseFilter) HandleWrite(ctx FilterContext) FilterAction {
	return ctx.GetNextAction()
}

func (self *BaseFilter) HandleOpen(ctx FilterContext) FilterAction {
	return ctx.GetNextAction()
}

func (self *BaseFilter) HandleClose(ctx FilterContext) FilterAction {
	return ctx.GetNextAction()
}

func (self *BaseFilter) HandleError(ctx FilterContext) FilterAction {
	return ctx.GetNextAction()
}
