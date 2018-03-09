package filter

import (
	"fairy"
	"fairy/base"
	"fairy/log"
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

func NewLogging() *LoggingFilter {
	return NewLoggingEx(LoggingFilterConnect | LoggingFilterClose | LoggingFilterError)
}

func NewLoggingEx(mask int) *LoggingFilter {
	lf := &LoggingFilter{}
	lf.mask = uint8(mask)
	return lf
}

type LoggingFilter struct {
	base.BaseFilter
	mask uint8
}

func (lf *LoggingFilter) need(m uint8) bool {
	return (lf.mask & m) != 0
}

func (lf *LoggingFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	if lf.need(LoggingFilterRead) {
		buffer, ok := ctx.GetMessage().(*fairy.Buffer)
		if ok {
			log.Debug("read data:len=%+v", buffer.Length())
		}
	}

	return ctx.GetNextAction()
}

func (lf *LoggingFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	if lf.need(LoggingFilterWrite) {
		buffer, ok := ctx.GetMessage().(*fairy.Buffer)
		if ok {
			log.Debug("send data:len=%+v", buffer.Length())
		}
	}

	return ctx.GetNextAction()
}

func (lf *LoggingFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConn()

	if lf.need(LoggingFilterAccept) && conn.IsServerSide() {
		log.Debug("accept new conn:id=%+v, addr=%+v", conn.GetConnId(), conn.RemoteAddr())
	}

	if lf.need(LoggingFilterConnect) && conn.IsClientSide() {
		log.Debug("connect success:id=%+v, addr=%+v", conn.GetConnId(), conn.RemoteAddr())
	}

	return ctx.GetNextAction()
}

func (lf *LoggingFilter) HandleClose(ctx fairy.FilterContext) fairy.FilterAction {
	if lf.need(LoggingFilterClose) {
		conn := ctx.GetConn()
		log.Debug("close conn:id=%+v, isclient=%+v", conn.GetConnId(), conn.IsClientSide())
	}

	return ctx.GetNextAction()
}

func (lf *LoggingFilter) HandleError(ctx fairy.FilterContext) fairy.FilterAction {
	if lf.need(LoggingFilterClose) {
		conn := ctx.GetConn()
		log.Error("conn error:id=%+v err=%+v", conn.GetConnId(), ctx.GetError())
	}

	return ctx.GetNextAction()
}
