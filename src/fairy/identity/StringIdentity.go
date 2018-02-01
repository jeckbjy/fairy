package identity

import (
	"errors"
	"fairy"
	"fairy/packet"
)

func NewStringIdentity() *StringIdentity {
	id := &StringIdentity{}
	return id
}

/**
 * StringIdentity 冒号分隔消息头和消息体
 */
type StringIdentity struct {
}

func (self *StringIdentity) Decode(buffer *fairy.Buffer) (fairy.Packet, error) {
	name, err := buffer.ReadUntil(':')
	if err != nil {
		return nil, err
	}

	pkt := packet.NewBasePacket()
	pkt.SetName(string(name))
	return pkt, nil
}

func (self *StringIdentity) Encode(buffer *fairy.Buffer, data interface{}) error {
	pkt, ok := data.(fairy.Packet)
	if !ok {
		return errors.New("StringString encode must packet!")
	}

	name := pkt.GetName()
	if name == "" {
		return errors.New("StringIdentity:cannot find name!")
	}

	buffer.Append([]byte(name))
	buffer.Append([]byte(":"))

	return nil
}
