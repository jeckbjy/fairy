package rpc

import (
	"sync"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/util"
)

var gRPCMgr *zRpcMgr

func init() {
	util.Once(gRPCMgr, func() {
		gRPCMgr = &zRpcMgr{}
	})
}

type zRpcMgr struct {
	handlers map[uint64]fairy.Handler
	mux      sync.Mutex
}

func (rm *zRpcMgr) Push(rpcid uint64, rh fairy.Handler) {
	rm.mux.Lock()
	rm.handlers[rpcid] = rh
	rm.mux.Unlock()
}

func (rm *zRpcMgr) Pop(rpcid uint64) fairy.Handler {
	var handler fairy.Handler
	rm.mux.Lock()
	if h, ok := rm.handlers[rpcid]; ok {
		handler = h
		delete(rm.handlers, rpcid)
	}
	rm.mux.Unlock()
	return handler
}
