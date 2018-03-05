package codec

import "fairy"

func NewBinary() *BinaryCodec {
	bc := &BinaryCodec{}
	return bc
}

type BinaryCodec struct {
}

func (bc *BinaryCodec) Encode(msg interface{}, buffer *fairy.Buffer) error {
	// data, err := json.Marshal(msg)
	// if err != nil {
	// 	return err
	// }

	// buffer.Append(data)
	return nil
}

func (bc *BinaryCodec) Decode(msg interface{}, buffer *fairy.Buffer) error {
	// data := buffer.ReadToEnd()
	// return json.Unmarshal(data, msg)
	return nil
}
