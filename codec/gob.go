package codec

import (
	"encoding/gob"

	"github.com/jeckbjy/fairy"
)

func NewGob() *GobCodec {
	codec := &GobCodec{}
	return codec
}

type GobCodec struct {
}

func (self *GobCodec) Encode(msg interface{}, buffer *fairy.Buffer) error {
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(msg)
	return err
}

func (self *GobCodec) Decode(msg interface{}, buffer *fairy.Buffer) error {
	dec := gob.NewDecoder(buffer)
	err := dec.Decode(msg)
	return err
}
