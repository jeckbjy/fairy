package filter

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

// NewConnect 连接成功回调
func NewConnect(cb ConnectCB) fairy.Filter {
	filter := &zConnectFilter{cb: cb, sync: true}
	return filter
}

// NewAsyncConnect 连接成功直接回调,不需要post主线程
func NewAsyncConnect(cb ConnectCB) fairy.Filter {
	filter := &zConnectFilter{cb: cb, sync: false}
	return filter
}

type ConnectCB func(fairy.Conn)

type zConnectFilter struct {
	base.BaseFilter
	cb   ConnectCB
	sync bool
}

func (self *zConnectFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConn()
	if conn.IsClientSide() && self.cb != nil {
		if self.sync {
			fairy.GetExecutor().Dispatch(fairy.NewFuncEvent(func() {
				self.cb(conn)
			}))
		} else {
			self.cb(conn)
		}
	}
	return ctx.GetNextAction()
}
