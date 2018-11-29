package identities

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/packet"
)

func NewFixed16() *Fixed16Identity {
	return &Fixed16Identity{}
}

// Fixed16Identity uint16保存消息ID,小端编码,0<id<65535
type Fixed16Identity struct {
}

func (*Fixed16Identity) Encode(buffer *fairy.Buffer, pkt fairy.IPacket) error {
	id := pkt.GetId()
	if id == 0 {
		return fmt.Errorf("packet id cannot zero")
	}

	if id >= math.MaxUint16 {
		return fmt.Errorf("packet id overflow:[id=%+v]", id)
	}

	data := make([]byte, 2)
	binary.LittleEndian.PutUint16(data, uint16(id))
	buffer.Append(data)
	return nil
}

func (*Fixed16Identity) Decode(buffer *fairy.Buffer) (fairy.IPacket, error) {
	// 读取ID
	data := make([]byte, 2)
	if _, err := io.ReadFull(buffer, data); err != nil {
		return nil, err
	}

	// 不要删除数据,将来可能还要用
	// buffer.Discard()

	// 解析ID
	id := binary.LittleEndian.Uint16(data)
	if id == 0 {
		return nil, fmt.Errorf("packet id cannot zero")
	}

	// create packet
	pkt := packet.NewBase()
	pkt.SetId(uint(id))

	return pkt, nil
}
