package identity

import (
	"fairy"
	"fairy/packet"
	"fairy/util"
	"fmt"
	"io"
	"math"
)

func NewInteger() *IntegerIdentity {
	return NewIntegerIdentityEx(true)
}

func NewIntegerIdentityEx(littleEndian bool) *IntegerIdentity {
	identity := &IntegerIdentity{}
	identity.littleEndian = littleEndian
	return identity
}

/**
 * 使用uint16保存消息ID，区分大小端，要求消息ID不能为0
 */
type IntegerIdentity struct {
	littleEndian bool
}

func (self *IntegerIdentity) Decode(buffer *fairy.Buffer) (fairy.Packet, error) {
	// 读取ID
	data := make([]byte, 2)
	if _, err := io.ReadFull(buffer, data); err != nil {
		return nil, err
	}

	// 不要删除数据,将来可能还要用
	// buffer.Discard()

	// 解析ID
	id := uint(util.GetUint16(data, self.littleEndian))
	if id == 0 {
		return nil, fmt.Errorf("IntegerIdentity:msgid is zero!")
	}

	// create packet
	pkt := packet.NewBase()
	pkt.SetId(id)

	return pkt, nil
}

func (self *IntegerIdentity) Encode(buffer *fairy.Buffer, data interface{}) error {
	pkt, ok := data.(fairy.Packet)
	if !ok {
		return fmt.Errorf("IntegerIdentity encode must be packet")
	}

	id := pkt.GetId()
	if id == 0 {
		return fmt.Errorf("IntegerIdentity encode id cannot zero!")
	}

	if id >= math.MaxUint16 {
		return fmt.Errorf("IntegerIdentity encode id[%v] overflow!", id)
	}

	buff := make([]byte, 2)
	util.PutUint16(buff, uint16(id), self.littleEndian)
	buffer.Append(buff)

	return nil
}
