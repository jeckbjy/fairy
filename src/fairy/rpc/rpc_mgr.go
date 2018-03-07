package rpc

import (
	"fairy"
	"fairy/util"
	"sync"
)

var gRpcMgr *RpcMgr

func init() {
	util.Once(gRpcMgr, func() {
		gRpcMgr = &RpcMgr{}
	})
}

type RpcMgr struct {
	handlers map[uint64]fairy.Handler
	mux      sync.Mutex
}

func (rm *RpcMgr) Push(rpcid uint64, rh *RpcHandler) {
	rm.mux.Lock()
	rm.handlers[rpcid] = rh
	rm.mux.Unlock()
}

func (rm *RpcMgr) Pop(rpcid uint64) fairy.Handler {
	var handler fairy.Handler
	rm.mux.Lock()
	if h, ok := rm.handlers[rpcid]; ok {
		handler = h
		delete(rm.handlers, rpcid)
	}
	rm.mux.Unlock()
	return handler
}
