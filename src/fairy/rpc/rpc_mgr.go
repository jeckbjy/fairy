package rpc

import (
	"fairy"
	"fairy/util"
	"sync"
)

var gRPCMgr *rpcMgr

func init() {
	util.Once(gRPCMgr, func() {
		gRPCMgr = &rpcMgr{}
	})
}

type rpcMgr struct {
	handlers map[uint64]fairy.Handler
	mux      sync.Mutex
}

func (rm *rpcMgr) Push(rpcid uint64, rh fairy.Handler) {
	rm.mux.Lock()
	rm.handlers[rpcid] = rh
	rm.mux.Unlock()
}

func (rm *rpcMgr) Pop(rpcid uint64) fairy.Handler {
	var handler fairy.Handler
	rm.mux.Lock()
	if h, ok := rm.handlers[rpcid]; ok {
		handler = h
		delete(rm.handlers, rpcid)
	}
	rm.mux.Unlock()
	return handler
}
