package frames

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/jeckbjy/fairy"
)

// NewVarint 创建变长消息编码
func NewVarint() *VarintFrame {
	frame := &VarintFrame{}
	return frame
}

// VarintFrame 消息长度使用变长编码
type VarintFrame struct {
}

func (*VarintFrame) Encode(buffer *fairy.Buffer) error {
	size := buffer.Length()
	data := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(data, uint64(size))
	buffer.Prepend(data[:n])
	return nil
}

func (*VarintFrame) Decode(buffer *fairy.Buffer) (*fairy.Buffer, error) {
	size, err := binary.ReadUvarint(buffer)
	if err != nil {
		return nil, err
	}

	// check data size
	if !buffer.HasRemain(int(size)) {
		return nil, fmt.Errorf("VarintLengthFrame no enouth data")
	}

	result := fairy.NewBuffer()

	// remove length head
	buffer.Discard()
	// read data
	buffer.Seek(int(size), io.SeekStart)
	buffer.Split(result)
	return result, nil
}
