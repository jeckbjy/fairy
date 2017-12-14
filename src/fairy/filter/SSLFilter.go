package filter

import (
	"fairy"
	"fairy/base"
)

type SSLFilter struct {
	base.BaseFilter
}

func (self *SSLFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	return ctx.GetNextAction()
}

func (self *SSLFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	return ctx.GetNextAction()
}
