package filter

import (
	"fairy"
	"fairy/base"
)

func NewFrameFilter(frame fairy.Frame) *FrameFilter {
	filter := &FrameFilter{}
	filter.Frame = frame
	return filter
}

type FrameFilter struct {
	base.BaseFilter
	fairy.Frame
}

func (self *FrameFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
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

func (self *FrameFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	if buffer, ok := ctx.GetMessage().(*fairy.Buffer); ok {
		err := self.Encode(buffer)
		if err != nil {
			return ctx.GetStopAction()
		}

		// ctx.SetMessage(msg)
	}

	return ctx.GetNextAction()
}
