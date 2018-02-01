package packet

import (
	"fairy"
)

func NewServer() *ServerPacket {
	pkt := &ServerPacket{}
	return pkt
}

// 服务器内部通信用Packet
type ServerPacket struct {
	NormalPacket
	Host string
	Uid  uint64
	Mode uint
}

func (self *ServerPacket) Encode(buffer *fairy.Buffer) error {
	return nil
}

func (self *ServerPacket) Decode(buffer *fairy.Buffer) error {
	return nil
}
