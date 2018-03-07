package filter

import (
	"fairy"
	"fairy/base"
)

func NewConnect(cb ConnectCallback) *ConnectFilter {
	filter := &ConnectFilter{cb: cb}
	return filter
}

type ConnectCallback func(fairy.Conn)

type ConnectFilter struct {
	base.BaseFilter
	cb ConnectCallback
}

func (self *ConnectFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn.IsClientSide() && self.cb != nil {
		self.cb(conn)
	}
	return ctx.GetNextAction()
}

func (self *ConnectFilter) HandleError(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn.IsClientSide() {
		// conn.Reconnect()
	}

	return ctx.GetNextAction()
}

func (self *ConnectFilter) HandleClose(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn.IsClientSide() {
		// conn.Reconnect()
	}
	return ctx.GetNextAction()
}
