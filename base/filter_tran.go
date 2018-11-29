package base

import (
	"github.com/jeckbjy/fairy"
)

func NewTransferFilter() *TransferFilter {
	return &TransferFilter{}
}

// TransferFilter 底层的消息发送与接收,必须放在Chain的第一个
type TransferFilter struct {
	Filter
}

func (tf *TransferFilter) Name() string {
	return "TransferFilter"
}

func (self *TransferFilter) HandleRead(ctx fairy.IFilterCtx) {
	// 获得读缓存
	conn := ctx.GetConn()
	data := conn.Read()
	if data == nil || data.Empty() {
		return
	}

	// 移动到开始位置
	// data.Seek(0, io.SeekStart)
	data.Rewind()
	ctx.SetData(data)
	ctx.Next()
}

func (self *TransferFilter) HandleWrite(ctx fairy.IFilterCtx) {
	// 底层发送
	if buffer, ok := ctx.GetData().(*fairy.Buffer); ok {
		ctx.GetConn().Write(buffer)
	}

	ctx.Next()
}
