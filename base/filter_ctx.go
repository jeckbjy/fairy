package base

import (
	"errors"
	"sync"

	"github.com/jeckbjy/fairy"
)

var errBadIndex = errors.New("filter index out of range")
var ctxPool = &FilterCtxPool{}

// FilterCtxPool context pool
type FilterCtxPool struct {
	items *FilterCtx
	count int
	mutex sync.Mutex
}

func (pool *FilterCtxPool) Alloc(chain fairy.IFilterChain, conn fairy.IConn, cb callback) *FilterCtx {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()

	var ctx *FilterCtx
	if pool.items != nil {
		ctx = pool.items
		pool.items = ctx.next
		ctx.next = nil
	} else {
		ctx = &FilterCtx{}
	}

	ctx.init(chain, conn, cb)
	return ctx
}

func (pool *FilterCtxPool) Free(ctx *FilterCtx) {
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	// TODO:auto shrink when count > 2 * connSize
	ctx.next = pool.items
	pool.items = ctx
	pool.count++
}

type callback func(ctx fairy.IFilterCtx, index int)

type FilterCtx struct {
	AttrMap
	chain fairy.IFilterChain
	cb    callback
	index int
	conn  fairy.IConn
	data  interface{}
	next  *FilterCtx
}

func (ctx *FilterCtx) init(chain fairy.IFilterChain, conn fairy.IConn, cb callback) {
	ctx.chain = chain
	ctx.cb = cb
	ctx.index = -1
	ctx.conn = conn
	ctx.data = nil
}

func (ctx *FilterCtx) GetConn() fairy.IConn {
	return ctx.conn
}

func (ctx *FilterCtx) SetData(data interface{}) {
	ctx.data = data
}

func (ctx *FilterCtx) GetData() interface{} {
	return ctx.data
}

// Error
func (ctx *FilterCtx) Error(err error) {
	ctx.chain.HandleError(ctx.conn, err)
}

func (ctx *FilterCtx) Next() {
	ctx.index++
	if ctx.index < ctx.chain.Len() {
		ctx.cb(ctx, ctx.index)
	}
}

func (ctx *FilterCtx) Jump(index int) error {
	ctx.index = index
	if ctx.index > -1 && ctx.index < ctx.chain.Len() {
		ctx.cb(ctx, ctx.index)
		return nil
	}

	return errBadIndex
}

func (ctx *FilterCtx) JumpBy(name string) error {
	return ctx.Jump(ctx.chain.IndexOf(name))
}
