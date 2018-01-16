package codecs

import (
	"fairy"

	"github.com/golang/protobuf/proto"
)

func NewProtobufCodec() fairy.Codec {
	codec := &ProtobufCodec{}
	return codec
}

type ProtobufCodec struct {
}

func (self *ProtobufCodec) Encode(obj interface{}, buffer *fairy.Buffer) error {
	msg := obj.(proto.Message)
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	buffer.Append(data)
	return nil
}

func (self *ProtobufCodec) Decode(obj interface{}, buffer *fairy.Buffer) error {
	data := buffer.ToBytes()
	msg := obj.(proto.Message)
	return proto.Unmarshal(data, msg)
}
