package rpc

import (
	"fairy"
	"fairy/base"
	"fairy/timer"
	"fairy/util"
	"fmt"
)

// PopHandler get and remove handler
func PopHandler(rpcid uint64) fairy.Handler {
	return gRPCMgr.Pop(rpcid)
}

// Call Remote procedure call,result future can sync
func Call(conn fairy.Conn, pkt fairy.Packet, timeout int64, cb fairy.HandlerCB) (fairy.Future, error) {
	// 必须大于0，否则的话应该注册到Dispatcher中效率更高
	if timeout <= 0 {
		return nil, fmt.Errorf("rpm timeout must be greater than zero")
	}

	// var rpcid int
	rpcid, err := util.NextID()
	if err != nil {
		return nil, err
	}

	promise := base.NewPromise(conn)
	pkt.SetRpcId(rpcid)
	handler := newHandler(cb, promise)
	gRPCMgr.Push(rpcid, handler)

	timer.Start(timeout, func(t *timer.Timer) {
		// checkt timeout
		rh := gRPCMgr.Pop(rpcid)
		if rh != nil {
			pkt.SetTimeout()
			rh.Invoke(conn, pkt)
		}
	})

	return promise, nil
}
