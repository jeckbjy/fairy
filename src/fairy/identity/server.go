package identity

import (
	"fairy"
	"fairy/packet"
	"fmt"
)

func NewServer() *ServerIdentity {
	return &ServerIdentity{}
}

// for server packet
type ServerIdentity struct {
}

func (self *ServerIdentity) Encode(buffer *fairy.Buffer, data interface{}) error {
	switch data.(type) {
	case *packet.BasePacket:
		return packet.EncodeNormal(data.(*packet.BasePacket), buffer)
	case *packet.ServerPacket:
		return packet.EncodeServer(data.(*packet.ServerPacket), buffer)
	default:
		return fmt.Errorf("encode must be ServerPacket")
	}

	return fmt.Errorf("encode must be NormalPacket")
}

func (self *ServerIdentity) Decode(buffer *fairy.Buffer) (fairy.Packet, error) {
	pkt := packet.NewServer()
	if err := packet.DecodeServer(pkt, buffer); err != nil {
		return nil, err
	}

	return pkt, nil
}
