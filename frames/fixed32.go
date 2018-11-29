package frames

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/jeckbjy/fairy"
)

func NewFixed32() *Fixed32Frame {
	return &Fixed32Frame{}
}

// Fixed32Frame 使用uint16保存,小端编码
type Fixed32Frame struct {
}

// Encode 头部追加消息体长度
func (*Fixed32Frame) Encode(buffer *fairy.Buffer) error {
	length := buffer.Length()
	if uint(length) >= math.MaxUint32 {
		return fmt.Errorf("frame length overflow:%+v", length)
	}
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(length))
	buffer.Prepend(data)
	return nil
}

//
func (*Fixed32Frame) Decode(buffer *fairy.Buffer) (*fairy.Buffer, error) {
	var data [4]byte
	if _, err := buffer.Read(data[:]); err != nil {
		return nil, fmt.Errorf("read head fail")
	}

	var length = binary.LittleEndian.Uint32(data[:])
	if !buffer.HasRemain(int(length)) {
		return nil, fmt.Errorf("no enough data")
	}

	result := fairy.NewBuffer()

	// discard length
	buffer.Discard()
	// read data
	buffer.Seek(int(length), io.SeekStart)
	buffer.Split(result)
	return result, nil
}
