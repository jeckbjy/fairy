package filter

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

// NewClose return CloseFilter
func NewClose(cb CloseCB) fairy.Filter {
	filter := &zCloseFilter{cb: cb, sync: true}
	return filter
}

func NewAsyncClose(cb CloseCB) fairy.Filter {
	filter := &zCloseFilter{cb: cb, sync: false}
	return filter
}

// CloseCB call when conn closed
type CloseCB func(fairy.Conn)

type zCloseFilter struct {
	base.BaseFilter
	cb   CloseCB
	sync bool
}

func (self *zCloseFilter) HandleClose(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConn()
	if self.cb != nil {
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
