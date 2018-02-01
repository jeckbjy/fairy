package identity

import (
	"fairy"
	"fairy/packet"
	"fmt"
)

// for server packet
type ServerIdentity struct {
}

func (self *ServerIdentity) Decode(buffer *fairy.Buffer) (fairy.Packet, error) {
	pkt := packet.NewServer()
	if err := pkt.Decode(buffer); err != nil {
		return nil, err
	}

	return nil, nil
}

func (self *ServerIdentity) Encode(buffer *fairy.Buffer, data interface{}) error {
	if pkt, ok := data.(*packet.NormalPacket); ok {
		return pkt.Encode(buffer)
	}

	// integer
	// if pkt, ok := data.(*packet.BasePacket); ok {

	// }

	return fmt.Errorf("encode must be NormalPacket")
}
