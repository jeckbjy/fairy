package filter

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

func NewFrame(frame fairy.Frame) *zFrameFilter {
	filter := &zFrameFilter{}
	filter.Frame = frame
	return filter
}

type zFrameFilter struct {
	base.BaseFilter
	fairy.Frame
}

func (self *zFrameFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	if buffer, ok := ctx.GetMessage().(*fairy.Buffer); ok {
		result, err := self.Decode(buffer)
		if err != nil {
			return ctx.GetStopAction()
		}

		// 透传buffer
		ctx.SetMessage(result)
	}

	return ctx.GetNextAction()
}

func (self *zFrameFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	if buffer, ok := ctx.GetMessage().(*fairy.Buffer); ok {
		err := self.Encode(buffer)
		if err != nil {
			return ctx.GetStopAction()
		}

		// ctx.SetMessage(msg)
	}

	return ctx.GetNextAction()
}
