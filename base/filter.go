package base

import (
	"github.com/jeckbjy/fairy"
)

// Filter 实现IFilter接口
type Filter struct {
}

func (f *Filter) HandleRead(ctx fairy.IFilterCtx) {
	ctx.Next()
}

func (f *Filter) HandleWrite(ctx fairy.IFilterCtx) {
	ctx.Next()
}

func (f *Filter) HandleOpen(ctx fairy.IFilterCtx) {
	ctx.Next()
}

func (f *Filter) HandleClose(ctx fairy.IFilterCtx) {
	ctx.Next()
}

func (f *Filter) HandleError(ctx fairy.IFilterCtx) {
	ctx.Next()
}
