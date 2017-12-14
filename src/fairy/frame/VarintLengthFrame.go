package frame

import (
	"encoding/binary"
	"errors"
	"fairy"
)

func NewVarintLengthFrame() *VarintLengthFrame {
	frame := &VarintLengthFrame{}
	return frame
}

type VarintLengthFrame struct {
}

func (self *VarintLengthFrame) Encode(buffer *fairy.Buffer) error {
	size := buffer.Length()
	data := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(data, uint64(size))
	buffer.Prepend(data[:n])
	return nil
}

func (self *VarintLengthFrame) Decode(buffer *fairy.Buffer) (*fairy.Buffer, error) {
	size, err := binary.ReadUvarint(buffer)
	if err != nil {
		return nil, err
	}

	// check data size
	if !buffer.HasRemain(int(size)) {
		return nil, errors.New("VarintLengthFrame no enouth data!")
	}

	result := fairy.NewBuffer()

	buffer.Discard()
	buffer.Split(result)
	return result, nil
}
