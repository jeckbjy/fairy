package base

import (
	"github.com/jeckbjy/fairy"
)

func NewChain() *FilterChain {
	chain := &FilterChain{}
	return chain
}

// FilterChain 实现fairy.FilterChain接口
type FilterChain struct {
	filters []fairy.IFilter
}

func (chain *FilterChain) Len() int {
	return len(chain.filters)
}

func (chain *FilterChain) IndexOf(name string) int {
	for index, filter := range chain.filters {
		if filter.Name() == name {
			return index
		}
	}

	return -1
}

func (chain *FilterChain) AddFirst(filters ...fairy.IFilter) {
	chain.filters = append(filters, chain.filters...)
}

func (chain *FilterChain) AddLast(filters ...fairy.IFilter) {
	chain.filters = append(chain.filters, filters...)
}

func (chain *FilterChain) HandleOpen(conn fairy.IConn) {
	ctx := ctxPool.Alloc(chain, conn, chain.doOpen)
	ctx.Next()
}

func (chain *FilterChain) HandleClose(conn fairy.IConn) {
	ctx := ctxPool.Alloc(chain, conn, chain.doClose)
	ctx.Next()
}

func (chain *FilterChain) HandleRead(conn fairy.IConn) {
	ctx := ctxPool.Alloc(chain, conn, chain.doRead)
	ctx.Next()
}

func (chain *FilterChain) HandleWrite(conn fairy.IConn, msg interface{}) {
	ctx := ctxPool.Alloc(chain, conn, chain.doWrite)
	ctx.SetData(msg)
	ctx.Next()
}

func (chain *FilterChain) HandleError(conn fairy.IConn, err error) {
	ctx := ctxPool.Alloc(chain, conn, chain.doError)
	ctx.SetData(err)
	ctx.Next()
}

func (chain *FilterChain) doOpen(ctx fairy.IFilterCtx, index int) {
	chain.filters[index].HandleOpen(ctx)
}

func (chain *FilterChain) doClose(ctx fairy.IFilterCtx, index int) {
	chain.filters[index].HandleClose(ctx)
}

func (chain *FilterChain) doRead(ctx fairy.IFilterCtx, index int) {
	chain.filters[index].HandleRead(ctx)
}

func (chain *FilterChain) doWrite(ctx fairy.IFilterCtx, index int) {
	chain.filters[index].HandleWrite(ctx)
}

func (chain *FilterChain) doError(ctx fairy.IFilterCtx, index int) {
	chain.filters[index].HandleError(ctx)
}
