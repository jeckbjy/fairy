package rpc

import (
	"fmt"

	"github.com/jeckbjy/fairy/base"
	"github.com/jeckbjy/fairy/timer"
	"github.com/jeckbjy/fairy/util"

	"github.com/jeckbjy/fairy"
)

// PopHandler get and remove handler
func PopHandler(rpcid uint64) fairy.Handler {
	return gRPCMgr.Pop(rpcid)
}

// Call Remote procedure call,result future can sync
func Call(conn fairy.Conn, pkt fairy.Packet, timeout int64, sync bool, cb fairy.HandlerCB) error {
	if timeout <= 0 {
		return fmt.Errorf("rpm timeout must be greater than zero")
	}

	rpcid, err := util.NextID()
	if err != nil {
		return err
	}

	var promise fairy.Promise
	if !sync {
		promise = base.NewPromise(conn)
	}

	pkt.SetRpcId(rpcid)
	handler := &RPCHandler{cb: cb, promise: promise}
	gRPCMgr.Push(rpcid, handler)

	timer.Start(timer.ModeDelay, timeout, func() {
		rh := gRPCMgr.Pop(rpcid)
		if rh != nil {
			pkt.SetTimeout()
			ctx := fairy.NewHandlerCtx(conn, pkt, rh, fairy.GetDispatcher().Middlewares())
			ctx.Process()
		}
	})

	if promise != nil {
		promise.Wait(timeout)
	}

	return nil
}
