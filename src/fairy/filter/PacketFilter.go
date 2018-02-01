package filter

import (
	"fairy"
	"fairy/base"
	"fairy/packet"
)

func NewPacketFilter(identity fairy.Identity, codec fairy.Codec) *PacketFilter {
	filter := &PacketFilter{}
	filter.Identity = identity
	filter.Codec = codec
	return filter
}

// decode:identity->UncaughtHandler->codec->dispatcher
type PacketFilter struct {
	base.BaseFilter
	fairy.Identity
	fairy.Codec
	*fairy.Registry
	*fairy.Dispatcher
}

func (self *PacketFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	// create buffer
	data := ctx.GetMessage()
	buffer, ok := data.(*fairy.Buffer)
	if !ok {
		return ctx.GetNextAction()
	}

	pkt, err := self.Identity.Decode(buffer)
	if err != nil {
		// throw error??
		return ctx.GetStopAction()
	}

	if pkt == nil {
		return ctx.GetNextAction()
	}

	// UncaughtHandler, donot need codec
	handler := self.Dispatcher.GetHandler(pkt.GetId(), pkt.GetName())
	if handler == nil {
		uncaught := self.Dispatcher.GetUncaughtHandler()
		ctx.SetHandler(uncaught)
		return ctx.GetNextAction()
	}

	// create msg
	msg := self.Registry.Create(pkt.GetId(), pkt.GetName())
	if msg == nil {
		return ctx.GetNextAction()
	}

	// codec
	err = self.Codec.Decode(msg, buffer)
	if err != nil {
		return ctx.GetStopAction()
	}

	ctx.SetMessage(pkt)
	ctx.SetHandler(handler)
	return ctx.GetNextAction()
}

func (self *PacketFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	data := ctx.GetMessage()
	if _, ok := data.(*fairy.Buffer); ok {
		// 已经是buffer，无需编码
		return ctx.GetNextAction()
	}

	buffer := fairy.NewBuffer()

	var pkt fairy.Packet
	var msg interface{}
	var ok bool

	pkt, ok = data.(fairy.Packet)
	if ok {
		msg = pkt.GetMessage()
	} else {
		id, name := self.Registry.GetInfo(data)
		pkt = &packet.BasePacket{}
		pkt.SetId(id)
		pkt.SetName(name)
		msg = data
	}

	// 写入头信息
	if err := self.Identity.Encode(buffer, pkt); err != nil {
		// throw error
		return ctx.GetStopAction()
	}

	// 写入消息体
	if err := self.Codec.Encode(msg, buffer); err != nil {
		// throw error
		return ctx.GetStopAction()
	}

	// 透传buffer
	ctx.SetMessage(buffer)

	return ctx.GetNextAction()
}
