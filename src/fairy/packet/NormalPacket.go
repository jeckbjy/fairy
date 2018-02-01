package packet

import (
	"fairy"
)

func NewNormal() *NormalPacket {
	pkt := &NormalPacket{}
	return pkt
}

type NormalPacket struct {
	BasePacket
}

func (self *NormalPacket) Encode(buffer *fairy.Buffer) error {
	return nil
}

func (self *NormalPacket) Decode(buffer *fairy.Buffer) error {
	return nil
}

// func EncodeNormalPacket(buffer *fairy.Buffer, packet *base.BasePacket) {
// 	codec := Codec{}
// 	codec.CreateReader(buffer)
// 	packet.SetResult(codec.ReadUInt())
// 	packet.SetSerialId(codec.ReadUInt())
// 	packet.SetId(codec.ReadUInt())
// }

// func DecodeNormalPacket(buffer *fairy.Buffer, packet *base.BasePacket) {
// 	codec := Codec{}
// 	codec.CreateWriter(buffer)
// 	codec.PushUInt(packet.GetResult())
// 	codec.PushUInt(packet.GetSerialId())
// 	codec.PushUInt(packet.GetId())
// }
