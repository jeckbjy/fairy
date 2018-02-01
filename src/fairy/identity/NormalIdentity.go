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

func (self *NormalIdentity) Decode(buffer *fairy.Buffer) (fairy.Packet, error) {
	pkt := packet.NewNormal()
	if err := pkt.Decode(buffer); err != nil {
		return nil, err
	}

	return nil, nil
}

func (self *NormalIdentity) Encode(buffer *fairy.Buffer, data interface{}) error {
	if pkt, ok := data.(*packet.NormalPacket); ok {
		return pkt.Encode(buffer)
	}

	// integer
	// if pkt, ok := data.(*packet.BasePacket); ok {

	// }

	return fmt.Errorf("encode must be NormalPacket")
}
