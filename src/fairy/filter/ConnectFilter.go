package filter

import (
	"fairy"
	"fairy/base"
)

func NewConnectFilter(cb ConnectCallback) *ConnectFilter {
	filter := &ConnectFilter{cb: cb}
	return filter
}

type ConnectCallback func(fairy.Connection)

type ConnectFilter struct {
	base.BaseFilter
	cb ConnectCallback
}

func (self *ConnectFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConnection()
	if conn.IsClientSide() {
		self.cb(conn)
	}
	return ctx.GetNextAction()
}
