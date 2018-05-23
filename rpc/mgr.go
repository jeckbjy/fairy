package rpc

import (
	"sync"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/util"
)

var gRPCMgr *RpcMgr

func init() {
	util.Once(gRPCMgr, func() {
		gRPCMgr = &RpcMgr{}
	})
}

// RpcMgr 管理RPC调用
type RpcMgr struct {
	handlers map[uint64]fairy.Handler
	mux      sync.Mutex
}

// Push 插入一条记录
func (rm *RpcMgr) Push(rpcid uint64, rh fairy.Handler) {
	rm.mux.Lock()
	rm.handlers[rpcid] = rh
	rm.mux.Unlock()
}

// Pop 删除一条记录
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
