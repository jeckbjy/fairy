package filter

import (
	"fairy"
	"fairy/base"
	"fairy/log"
)

func NewLog() *LogFilter {
	return &LogFilter{}
}

type LogFilter struct {
	base.BaseFilter
}

func (self *LogFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	buffer, ok := ctx.GetMessage().(*fairy.Buffer)
	if ok {
		log.Debug("read data:len=%+v", buffer.Length())
	}
	return ctx.GetNextAction()
}

func (self *LogFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	buffer, ok := ctx.GetMessage().(*fairy.Buffer)
	if ok {
		log.Debug("send data:len=%+v", buffer.Length())
	}

	return ctx.GetNextAction()
}

func (self *LogFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn != nil {
		log.Debug("open conn:id=%+v,isclient=%+v", conn.GetConnId(), conn.IsClientSide())
	}
	return ctx.GetNextAction()
}

func (self *LogFilter) HandleClose(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn != nil {
		log.Debug("close conn:id=%+v, isclient=%+v", conn.GetConnId(), conn.IsClientSide())
	}

	return ctx.GetNextAction()
}

func (self *LogFilter) HandleError(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn != nil {
		log.Error("connid=%+v, error=%+v", conn.GetConnId(), ctx.GetError())
	} else {
		log.Error("%+v", ctx.GetError())
	}

	return ctx.GetNextAction()
}
