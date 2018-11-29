package filters

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

func NewFrame(frame fairy.IFrame) *FrameFilter {
	return &FrameFilter{frame: frame}
}

// FrameFilter 粘包处理
type FrameFilter struct {
	base.Filter
	frame fairy.IFrame
}

func (ff *FrameFilter) Name() string {
	return "FrameFilter"
}

func (ff *FrameFilter) HandleRead(ctx fairy.IFilterCtx) {
	if buffer, ok := ctx.GetData().(*fairy.Buffer); ok {
		result, err := ff.frame.Decode(buffer)
		if err != nil {
			return
		}

		ctx.SetData(result)
	}

	ctx.Next()
}

func (ff *FrameFilter) HandleWrite(ctx fairy.IFilterCtx) {
	if buffer, ok := ctx.GetData().(*fairy.Buffer); ok {
		err := ff.frame.Encode(buffer)
		if err != nil {
			return
		}
	}

	ctx.Next()
}
