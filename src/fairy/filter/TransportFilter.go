package filter

import (
	"fairy"
	"fairy/base"
)

func NewTransportFilter() *TransportFilter {
	filter := &TransportFilter{}
	return filter
}

type TransportFilter struct {
	base.BaseFilter
}

func (self *TransportFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	// 先获得buffer
	conn := ctx.GetConnection()
	buffer := conn.Read()
	if buffer == nil || buffer.Empty() {
		return ctx.GetStopAction()
	}

	buffer.Rewind()
	ctx.SetMessage(buffer)
	return ctx.GetNextAction()
}

func (self *TransportFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	// 底层发送
	if buffer, ok := ctx.GetMessage().(*fairy.Buffer); ok {
		conn := ctx.GetConnection()
		conn.Write(buffer)
	}

	return ctx.GetNextAction()
}
