package rpc

import (
	"github.com/jeckbjy/fairy"
)

type RPCHandler struct {
	cb      fairy.HandlerCB
	promise fairy.Promise
}

func (h *RPCHandler) QueueID() int {
	return 0
}

func (h *RPCHandler) Invoke(ctx *fairy.HandlerCtx) {
	h.cb(ctx)

	if h.promise != nil {
		if ctx.IsSuccess() {
			h.promise.SetSuccess()
		} else {
			h.promise.SetFailure()
		}
	}
}
