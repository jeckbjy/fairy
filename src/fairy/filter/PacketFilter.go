package filter

import (
	"fairy"
	"fairy/base"
)

func NewPacketFilter(identity fairy.Identity, codec fairy.Codec) *PacketFilter {
	filter := &PacketFilter{}
	filter.Identity = identity
	filter.Codec = codec
	return filter
}

type PacketFilter struct {
	base.BaseFilter
	fairy.Identity
	fairy.Codec
}

func (self *PacketFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	// create buffer
	data := ctx.GetMessage()
	if buffer, ok := data.(*fairy.Buffer); ok {
		packet, err := self.Identity.Decode(buffer)
		if err != nil {
			return ctx.GetStopAction()
		}

		err = self.Codec.Decode(packet.GetMessage(), buffer)
		if err != nil {
			return ctx.GetStopAction()
		}
		// 透传下一个执行
		ctx.SetMessage(packet)
	}

	return ctx.GetNextAction()
}

func (self *PacketFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	message := ctx.GetMessage()
	if _, ok := message.(*fairy.Buffer); ok {
		// 已经是buffer，无需编码
		return ctx.GetNextAction()
	}

	buffer := fairy.NewBuffer()

	// 先写入头部信息:内部需要支持两种形式，packet or message
	if err := self.Identity.Encode(buffer, message); err != nil {
		return ctx.GetStopAction()
	}

	// 写需要支持两种，packet或者message
	if packet, ok := message.(fairy.Packet); ok {
		message = packet.GetMessage()
	}

	// 写入消息体
	if err := self.Codec.Encode(message, buffer); err != nil {
		return ctx.GetStopAction()
	}

	// 透传buffer
	ctx.SetMessage(buffer)

	return ctx.GetNextAction()
}
