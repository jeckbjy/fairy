package identity

import (
	"fairy"
	"fairy/base"
	"fairy/util"
	"fmt"
	"io"
	"math"
)

func NewDefaultIntegerIdentity() *IntegerIdentity {
	return NewIntegerIdentity(fairy.GetRegistry(), true)
}

func NewIntegerIdentity(registry *fairy.Registry, littleEndian bool) *IntegerIdentity {
	identity := &IntegerIdentity{}
	identity.Registry = registry
	identity.littleEndian = littleEndian
	return identity
}

// 使用uint16保存消息ID，区分大小端，要求消息ID不能为0
// encode如果不是packet
type IntegerIdentity struct {
	*fairy.Registry
	littleEndian bool
}

func (self *IntegerIdentity) Decode(buffer *fairy.Buffer) (fairy.Packet, error) {
	// 读取ID
	data := make([]byte, 2)
	if _, err := io.ReadFull(buffer, data); err != nil {
		return nil, err
	}

	// 解析ID
	id := uint(util.GetUint16(data, self.littleEndian))
	if id == 0 {
		return nil, fmt.Errorf("IntegerIdentity:msgid is zero!")
	}

	// create packet
	packet := base.NewBasePacket()
	packet.SetId(id)

	// create message
	msg := self.CreateById(id)
	if msg == nil {
		return packet, fmt.Errorf("IntegerIdentity: create message fail,msgid=%v!", id)
	}
	packet.SetMessage(msg)
	return packet, nil
}

func (self *IntegerIdentity) Encode(buffer *fairy.Buffer, data interface{}) error {
	id := uint(0)
	if packet, ok := data.(fairy.Packet); ok {
		id = packet.GetId()
	}

	if id == 0 {
		id = self.GetId(data)
	}

	if id == 0 {
		return fmt.Errorf("IntegerIdentity:cannot find msgid!")
	}

	if id >= math.MaxUint16 {
		return fmt.Errorf("IntegerIdentity:msgid overflow!")
	}

	id_buff := make([]byte, 2)
	util.PutUint16(id_buff, uint16(id), self.littleEndian)
	buffer.Append(id_buff)

	return nil
}
