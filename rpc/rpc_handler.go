package rpc

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/log"
)

func newHandler(cb fairy.HandlerCB, pro fairy.Promise) *zRpcHandler {
	rh := &zRpcHandler{cb: cb, promise: pro}
	return rh
}

type zRpcHandler struct {
	cb      fairy.HandlerCB
	promise fairy.Promise
}

func (rh *zRpcHandler) GetQueueId() int {
	return 0
}

func (rh *zRpcHandler) Invoke(conn fairy.Conn, packet fairy.Packet) {
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
