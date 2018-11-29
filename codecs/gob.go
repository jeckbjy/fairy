package codecs

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

func (self *GobCodec) Encode(buffer *fairy.Buffer, msg interface{}) error {
	enc := gob.NewEncoder(buffer)
	err := enc.Encode(msg)
	return err
}

func (self *GobCodec) Decode(buffer *fairy.Buffer, msg interface{}) error {
	dec := gob.NewDecoder(buffer)
	err := dec.Decode(msg)
	return err
}
