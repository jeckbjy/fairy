package filter

import (
	"fairy"
	"fairy/base"
)

func NewEncryptFilter() *EncryptFilter {
	filter := &EncryptFilter{}
	return filter
}

type EncryptFilter struct {
	base.BaseFilter
}

func (self *EncryptFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	return ctx.GetNextAction()
}

func (self *EncryptFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	return ctx.GetNextAction()
}
