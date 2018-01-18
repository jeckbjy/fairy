package packet

import (
	"fairy"
	"fairy/base"
)

// 服务器内部通信用Packet
type ServerPacket struct {
	base.BasePacket
	Host string
	Uid  uint64
	Mode uint
}

func EncodeServerPacket(buffer *fairy.Buffer, packet *ServerPacket) {

}

func DecodeServerPacket(buffer *fairy.Buffer, packet *ServerPacket) {

}
