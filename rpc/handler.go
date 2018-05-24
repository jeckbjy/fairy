package rpc

import (
	"github.com/jeckbjy/fairy"
)

type Handler struct {
	cb      fairy.HandlerCB
	promise fairy.Promise
}

func (h *Handler) QueueID() int {
	return 0
}

func (h *Handler) Invoke(ctx *fairy.HandlerCtx) {
	h.cb(ctx)

	if h.promise != nil {
		if ctx.IsSuccess() {
			h.promise.SetSuccess()
		} else {
			h.promise.SetFailure()
		}
	}
}
