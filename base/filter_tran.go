package base

import "github.com/jeckbjy/fairy"

func NewTransferFilter() *TransferFilter {
	filter := &TransferFilter{}
	return filter
}

type TransferFilter struct {
	BaseFilter
}

func (self *TransferFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	// 先获得buffer
	conn := ctx.GetConn()
	buffer := conn.Read()
	if buffer == nil || buffer.Empty() {
		return ctx.GetStopAction()
	}

	buffer.Rewind()
	ctx.SetMessage(buffer)
	return ctx.GetNextAction()
}

func (self *TransferFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	// 底层发送
	if buffer, ok := ctx.GetMessage().(*fairy.Buffer); ok {
		conn := ctx.GetConn()
		conn.Write(buffer)
	}

	return ctx.GetNextAction()
}
