package identity

import (
	"fmt"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/packet"
)

func NewString() *StringIdentity {
	id := &StringIdentity{}
	return id
}

/**
 * StringIdentity 冒号分隔消息头和消息体
 */
type StringIdentity struct {
}

func (self *StringIdentity) Encode(buffer *fairy.Buffer, data interface{}) error {
	pkt, ok := data.(fairy.Packet)
	if !ok {
		return fmt.Errorf("StringIdentity encode must be packet!")
	}

	name := pkt.GetName()
	if name == "" {
		return fmt.Errorf("StringIdentity:cannot find name!")
	}

	buffer.Append([]byte(name))
	buffer.Append([]byte(":"))

	return nil
}

func (self *StringIdentity) Decode(buffer *fairy.Buffer) (fairy.Packet, error) {
	name, err := buffer.ReadUntil(':')
	if err != nil {
		return nil, err
	}

	pkt := packet.NewBase()
	pkt.SetName(string(name))
	return pkt, nil
}
