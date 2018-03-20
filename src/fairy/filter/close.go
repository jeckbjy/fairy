package filter

import (
	"fairy"
	"fairy/base"
	"fairy/exec"
)

// NewClose return CloseFilter
func NewClose(cb CloseCB) fairy.Filter {
	filter := &CloseFilter{cb: cb, sync: true}
	return filter
}

func NewAsyncClose(cb CloseCB) fairy.Filter {
	filter := &CloseFilter{cb: cb, sync: false}
	return filter
}

// CloseCB call when conn closed
type CloseCB func(fairy.Conn)

type CloseFilter struct {
	base.BaseFilter
	cb   CloseCB
	sync bool
}

func (self *CloseFilter) HandleClose(ctx fairy.FilterContext) fairy.FilterAction {
	conn := ctx.GetConn()
	if self.cb != nil {
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
