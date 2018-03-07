package rpc

import (
	"fairy"
	"fairy/log"
)

func newHandler(cb fairy.HandlerCB, pro fairy.Promise) *RpcHandler {
	rh := &RpcHandler{cb: cb, promise: pro}
	return rh
}

type RpcHandler struct {
	cb      fairy.HandlerCB
	promise fairy.Promise
}

func (rh *RpcHandler) GetQueueId() int {
	return 0
}

func (rh *RpcHandler) Invoke(conn fairy.Conn, packet fairy.Packet) {
	defer log.Catch()
	rh.cb(conn, packet)

	if packet.IsSuccess() {
		rh.promise.SetSuccess()
	} else {
		rh.promise.SetFailure()
	}
}
