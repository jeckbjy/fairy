package filter

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
	"github.com/jeckbjy/fairy/packet"
	"github.com/jeckbjy/fairy/rpc"
)

func NewPacket(identity fairy.Identity, codec fairy.Codec) *zPacketFilter {
	filter := &zPacketFilter{}
	filter.Identity = identity
	filter.Codec = codec
	filter.Registry = fairy.GetRegistry()
	filter.Dispatcher = fairy.GetDispatcher()
	return filter
}

// decode:identity->UncaughtHandler->codec->dispatcher
type zPacketFilter struct {
	base.BaseFilter
	fairy.Identity
	fairy.Codec
	*fairy.Registry
	*fairy.Dispatcher
}

func (self *zPacketFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	// create buffer
	data := ctx.GetMessage()
	buffer, ok := data.(*fairy.Buffer)
	if !ok {
		return ctx.GetNextAction()
	}

	pkt, err := self.Identity.Decode(buffer)
	if err != nil {
		ctx.ThrowError(err)
		return ctx.GetStopAction()
	}

	if pkt == nil {
		return ctx.GetNextAction()
	}

	// find handler
	var handler fairy.Handler
	if pkt.GetRpcId() != 0 {
		handler = rpc.PopHandler(pkt.GetRpcId())
	}

	// UncaughtHandler, donot need codec
	if handler == nil {
		var ok bool
		handler, ok = self.Dispatcher.GetFinalHandler(pkt.GetId(), pkt.GetName())
		if !ok {
			ctx.SetHandler(handler)
			return ctx.GetNextAction()
		}
	}

	// create msg
	msg := self.Registry.Create(pkt.GetId(), pkt.GetName())
	if msg == nil {
		return ctx.GetNextAction()
	}

	// codec
	err = self.Codec.Decode(msg, buffer)
	if err != nil {
		ctx.ThrowError(err)
		return ctx.GetStopAction()
	}

	pkt.SetMessage(msg)

	ctx.SetMessage(pkt)
	ctx.SetHandler(handler)
	return ctx.GetNextAction()
}

func (self *zPacketFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	data := ctx.GetMessage()
	if _, ok := data.(*fairy.Buffer); ok {
		// 已经是buffer，无需编码
		return ctx.GetNextAction()
	}

	buffer := fairy.NewBuffer()

	var pkt fairy.Packet
	var msg interface{}
	var ok bool

	// 校验是packet,还是原始message
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
		ctx.ThrowError(err)
		return ctx.GetStopAction()
	}

	// 写入消息体
	if err := self.Codec.Encode(msg, buffer); err != nil {
		// throw error
		ctx.ThrowError(err)
		return ctx.GetStopAction()
	}

	// 透传buffer
	ctx.SetMessage(buffer)

	return ctx.GetNextAction()
}
