package rpc

import (
	"fairy"
	"fairy/log"
)

func newHandler(cb fairy.HandlerCB, pro fairy.Promise) *rpcHandler {
	rh := &rpcHandler{cb: cb, promise: pro}
	return rh
}

type rpcHandler struct {
	cb      fairy.HandlerCB
	promise fairy.Promise
}

func (rh *rpcHandler) GetQueueId() int {
	return 0
}

func (rh *rpcHandler) Invoke(conn fairy.Conn, packet fairy.Packet) {
	defer log.Catch()
	rh.cb(conn, packet)

	if rh.promise != nil {
		if packet.IsSuccess() {
			rh.promise.SetSuccess()
		} else {
			rh.promise.SetFailure()
		}
	}
}
