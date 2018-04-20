package filter

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
	"github.com/jeckbjy/fairy/log"
)

const (
	LoggingFilterRead    = 0x01
	LoggingFilterWrite   = 0x02
	LoggingFilterAccept  = 0x04
	LoggingFilterConnect = 0x08
	LoggingFilterClose   = 0x10
	LoggingFilterError   = 0x20
	LoggingFilterAll     = 0xff
)

func NewLogging() fairy.Filter {
	return NewLoggingEx(LoggingFilterConnect | LoggingFilterClose | LoggingFilterError)
}

func NewLoggingEx(mask int) fairy.Filter {
	lf := &zLoggingFilter{}
	lf.mask = uint8(mask)
	return lf
}

type zLoggingFilter struct {
	base.BaseFilter
	mask uint8
}

func (lf *zLoggingFilter) need(m uint8) bool {
	return (lf.mask & m) != 0
}

func (lf *zLoggingFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	if lf.need(LoggingFilterRead) {
		buffer, ok := ctx.GetMessage().(*fairy.Buffer)
		if ok {
			log.Debug("read data:len=%+v", buffer.Length())
		}
	}

	return ctx.GetNextAction()
}

func (lf *zLoggingFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	if lf.need(LoggingFilterWrite) {
		buffer, ok := ctx.GetMessage().(*fairy.Buffer)
		if ok {
			log.Debug("send data:len=%+v", buffer.Length())
		}
	}

	return ctx.GetNextAction()
}

func (lf *zLoggingFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConn()

	if lf.need(LoggingFilterAccept) && conn.IsServerSide() {
		log.Debug("accept new conn:id=%+v, addr=%+v", conn.GetConnId(), conn.RemoteAddr())
	}

	if lf.need(LoggingFilterConnect) && conn.IsClientSide() {
		log.Debug("connect success:id=%+v, addr=%+v", conn.GetConnId(), conn.RemoteAddr())
	}

	return ctx.GetNextAction()
}

func (lf *zLoggingFilter) HandleClose(ctx fairy.FilterContext) fairy.FilterAction {
	if lf.need(LoggingFilterClose) {
		conn := ctx.GetConn()
		log.Debug("close conn:id=%+v, isclient=%+v", conn.GetConnId(), conn.IsClientSide())
	}

	return ctx.GetNextAction()
}

func (lf *zLoggingFilter) HandleError(ctx fairy.FilterContext) fairy.FilterAction {
	if lf.need(LoggingFilterClose) {
		conn := ctx.GetConn()
		log.Error("conn error:id=%+v err=%+v", conn.GetConnId(), ctx.GetError())
	}

	return ctx.GetNextAction()
}
