package filter

import (
	"fairy"
	"fairy/base"
)

func NewLogFilter() *LogFilter {
	return &LogFilter{}
}

type LogFilter struct {
	base.BaseFilter
}

func (self *LogFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	buffer, ok := ctx.GetMessage().(*fairy.Buffer)
	if ok {
		fairy.Debug("read data:%+v", buffer.String())
	}
	return ctx.GetNextAction()
}

func (self *LogFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	buffer, ok := ctx.GetMessage().(*fairy.Buffer)
	if ok {
		fairy.Debug("send data:%+v", buffer.String())
	}

	return ctx.GetNextAction()
}

func (self *LogFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn != nil {
		fairy.Debug("open conn:id=%+v,isclient=%+v", conn.GetConnId(), conn.IsClientSide())
	}
	return ctx.GetNextAction()
}

func (self *LogFilter) HandleClose(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn != nil {
		fairy.Debug("close conn:id=%+v, isclient=%+v", conn.GetConnId(), conn.IsClientSide())
	}

	return ctx.GetNextAction()
}

func (self *LogFilter) HandleError(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn != nil {
		fairy.Error("connid=%+v, error=%+v", conn.GetConnId(), ctx.GetError())
	} else {
		fairy.Error("%+v", ctx.GetError())
	}

	return ctx.GetNextAction()
}
