package codec

import (
	"fairy"
	"fairy/codec/binary"
)

// NewBinary return BinaryCodec
func NewBinary() *BinaryCodec {
	bc := &BinaryCodec{}
	return bc
}

type BinaryCodec struct {
}

func (bc *BinaryCodec) Encode(msg interface{}, buffer *fairy.Buffer) error {
	data, err := binary.Marshal(msg)
	if err != nil {
		return err
	}

	buffer.Append(data)
	return nil
}

func (bc *BinaryCodec) Decode(msg interface{}, buffer *fairy.Buffer) error {
	data := buffer.ReadToEnd()
	return binary.Unmarshal(data, msg)
}
