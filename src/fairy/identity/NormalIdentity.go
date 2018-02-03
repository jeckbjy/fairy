package identity

import (
	"fairy"
	"fairy/packet"
	"fmt"
)

func NewNormal() *NormalIdentity {
	identity := &NormalIdentity{}
	return identity
}

// for normal packet
type NormalIdentity struct {
}

func (self *NormalIdentity) Encode(buffer *fairy.Buffer, data interface{}) error {
	if pkt, ok := data.(*packet.BasePacket); ok {
		return packet.EncodeNormal(pkt, buffer)
	}

	return fmt.Errorf("encode must be BasePacket")
}

func (self *NormalIdentity) Decode(buffer *fairy.Buffer) (fairy.Packet, error) {
	pkt := packet.NewNormal()
	if err := packet.DecodeNormal(pkt, buffer); err != nil {
		return nil, err
	}

	return pkt, nil
}
