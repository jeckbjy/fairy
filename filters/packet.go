package filters

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
	"github.com/jeckbjy/fairy/packet"
)

func NewPacket(identity fairy.IIdentity, codec fairy.ICodec) *PacketFilter {
	filter := &PacketFilter{}
	filter.identity = identity
	filter.codec = codec
	filter.registry = fairy.GetRegistry()
	return filter
}

// PacketFilter 创建并解析Packet
type PacketFilter struct {
	base.Filter
	identity fairy.IIdentity
	codec    fairy.ICodec
	registry *fairy.Registry
}

func (pf *PacketFilter) Name() string {
	return "PacketFilter"
}

// HandleRead 从Buffer中解析Packet
func (pf *PacketFilter) HandleRead(ctx fairy.IFilterCtx) {
	data := ctx.GetData()
	buffer, ok := data.(*fairy.Buffer)
	if !ok {
		ctx.Next()
		return
	}

	// 通过Identity创建Packet
	pkt, err := pf.identity.Decode(buffer)
	if err != nil || pkt == nil {
		// warning?
		return
	}

	// Tips:自定义Filter时,这里可以做一下优化
	// 如果没有定义Handler则不需要解析Message,仅仅用作转发处理

	// 解析消息包
	msg := pf.registry.Create(pkt.GetId(), pkt.GetName())
	if msg == nil {
		// warning?
		return
	}

	err = pf.codec.Decode(buffer, msg)
	if err != nil {
		// warning?
		return
	}

	pkt.SetMessage(msg)
	ctx.SetData(pkt)
	ctx.Next()
}

func (self *PacketFilter) HandleWrite(ctx fairy.IFilterCtx) {
	data := ctx.GetData()
	if _, ok := data.(*fairy.Buffer); ok {
		// 已经是buffer，无需编码
		ctx.Next()
		return
	}

	buffer := fairy.NewBuffer()

	var pkt fairy.IPacket
	var msg interface{}
	var ok bool

	// 校验是packet,还是原始message
	pkt, ok = data.(fairy.IPacket)
	if ok {
		msg = pkt.GetMessage()
	} else {
		// 原始消息,并非Packet,使用基础的Packet
		id, name, ok := self.registry.GetInfo(data)
		if !ok {
			// warning
			return
		}

		pkt = packet.NewBase()
		pkt.SetId(id)
		pkt.SetName(name)
		msg = data
	}

	// 写入头信息
	if err := self.identity.Encode(buffer, pkt); err != nil {
		// warning
		return
	}

	// 写入消息体
	if err := self.codec.Encode(buffer, msg); err != nil {
		// warning
		return
	}

	// 透传buffer
	ctx.SetData(buffer)
	ctx.Next()
}
