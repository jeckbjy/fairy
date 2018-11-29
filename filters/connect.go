package filters

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

// NewConnect 连接成功回调
func NewConnect(cb ConnectCB) *ConnectFilter {
	return &ConnectFilter{cb: cb}
}

// ConnectCB 连接回调函数
type ConnectCB func(fairy.IConn)

// ConnectFilter 连接成功时调用回调,注意:并非线程安全
type ConnectFilter struct {
	base.Filter
	cb ConnectCB
}

func (cf *ConnectFilter) Name() string {
	return "ConnectFilter"
}

func (cf *ConnectFilter) HandleOpen(ctx fairy.IFilterCtx) {
	conn := ctx.GetConn()
	if conn.IsConnector() {
		cf.cb(conn)
	}

	ctx.Next()
}
