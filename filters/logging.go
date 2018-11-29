package filters

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/log"
)

const (
	LoggingRead    = 0x01
	LoggingWrite   = 0x02
	LoggingAccept  = 0x04
	LoggingConnect = 0x08
	LoggingClose   = 0x10
	LoggingError   = 0x20
	LoggingAll     = 0xff
)

func NewLogging() *LoggingFilter {
	return &LoggingFilter{mask: LoggingAll}
}

type LoggingFilter struct {
	mask uint8
}

func (lf *LoggingFilter) need(m uint8) bool {
	return (lf.mask & m) != 0
}

func (lf *LoggingFilter) SetMask(m int) {
	lf.mask = uint8(m)
}

func (lf *LoggingFilter) Name() string {
	return "LoggingFilter"
}

func (lf *LoggingFilter) HandleRead(ctx fairy.IFilterCtx) {
	if lf.need(LoggingRead) {
		if buffer, ok := ctx.GetData().(*fairy.Buffer); ok {
			log.Debug("read data:len=%+v\n", buffer.Length())
		}
	}
	ctx.Next()
}

func (lf *LoggingFilter) HandleWrite(ctx fairy.IFilterCtx) {
	if lf.need(LoggingWrite) {
		if buffer, ok := ctx.GetData().(*fairy.Buffer); ok {
			log.Debug("send data:len=%+v", buffer.Length())
		}
	}

	ctx.Next()
}

func (lf *LoggingFilter) HandleOpen(ctx fairy.IFilterCtx) {
	conn := ctx.GetConn()

	if conn == nil {
		ctx.Next()
		return
	}

	if lf.need(LoggingAccept) && !conn.IsConnector() {
		log.Debug("accept new conn:id=%+v, addr=%+v", conn.GetId(), conn.RemoteAddr())
	}

	if lf.need(LoggingConnect) && conn.IsConnector() {
		log.Debug("connect success:id=%+v, addr=%+v", conn.GetId(), conn.RemoteAddr())
	}

	ctx.Next()
}

func (lf *LoggingFilter) HandleClose(ctx fairy.IFilterCtx) {
	conn := ctx.GetConn()
	if lf.need(LoggingClose) && conn != nil {
		log.Debug("close conn:id=%+v, isconnector=%+v", conn.GetId(), conn.IsConnector())
	}

	ctx.Next()
}

func (lf *LoggingFilter) HandleError(ctx fairy.IFilterCtx) {
	conn := ctx.GetConn()
	if lf.need(LoggingError) && conn != nil {
		log.Error("conn error:id=%+v err=%+v", conn.GetId(), ctx.GetData())
	}

	ctx.Next()
}
