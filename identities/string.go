package identities

import (
	"fmt"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/packet"
)

func NewString() *StringIdentity {
	return &StringIdentity{}
}

// StringIdentity 冒号分隔消息头和消息体,要求Packet中不能存在冒号
type StringIdentity struct {
}

func (*StringIdentity) Encode(buffer *fairy.Buffer, pkt fairy.IPacket) error {
	name := pkt.GetName()
	if name == "" {
		return fmt.Errorf("packet name is empty")
	}

	buffer.Append([]byte(name))
	buffer.Append([]byte(":"))
	return nil
}

func (*StringIdentity) Decode(buffer *fairy.Buffer) (fairy.IPacket, error) {
	name, err := buffer.ReadUntil(':')
	if err != nil {
		return nil, err
	}

	pkt := packet.NewBase()
	pkt.SetName(name)
	return pkt, nil
}
