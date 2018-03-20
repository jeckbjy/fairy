package filter

import (
	"fairy"
	"fairy/base"
	"fairy/exec"
)

// NewConnect 连接成功回调
func NewConnect(cb ConnectCB) fairy.Filter {
	filter := &ConnectFilter{cb: cb, sync: true}
	return filter
}

// NewAsyncConnect 连接成功直接回调,不需要post主线程
func NewAsyncConnect(cb ConnectCB) fairy.Filter {
	filter := &ConnectFilter{cb: cb, sync: false}
	return filter
}

type ConnectCB func(fairy.Conn)

type ConnectFilter struct {
	base.BaseFilter
	cb   ConnectCB
	sync bool
}

func (self *ConnectFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConn()
	if conn.IsClientSide() && self.cb != nil {
		if self.sync {
			exec.GetExecutor().Dispatch(exec.NewFuncEvent(func() {
				self.cb(conn)
			}))
		} else {
			self.cb(conn)
		}
	}
	return ctx.GetNextAction()
}
