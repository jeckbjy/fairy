package identity

import (
	"errors"
	"fairy"
	"fairy/base"
)

func NewDefaultStringIdentity() *StringIdentity {
	return NewStringIdentity(fairy.GetGlobalRegistry())
}

func NewStringIdentity(registry *fairy.Registry) *StringIdentity {
	identity := &StringIdentity{}
	identity.Registry = registry
	return identity
}

// name:body
type StringIdentity struct {
	*fairy.Registry
}

func (self *StringIdentity) Decode(buffer *fairy.Buffer) (fairy.Packet, error) {
	pos := buffer.IndexOf(':')
	if pos == -1 {
		return nil, nil
	}

	// read name
	name := ""

	// create packet
	packet := base.NewBasePacket()
	packet.SetName(name)

	// create message
	msg := self.CreateByName(name)
	if msg == nil {
		return packet, nil
	}

	packet.SetMessage(msg)
	return packet, nil
}

func (self *StringIdentity) Encode(buffer *fairy.Buffer, data interface{}) error {
	name := ""
	if packet, ok := data.(fairy.Packet); ok {
		name = packet.GetName()
	} else {
		name = self.GetName(data)
	}

	if name == "" {
		return errors.New("StringIdentity:cannot find name!")
	}

	buffer.Append([]byte(name))
	buffer.Append([]byte(":"))

	return nil
}
