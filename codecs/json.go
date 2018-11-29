package codecs

import (
	"encoding/json"

	"github.com/jeckbjy/fairy"
)

func NewJson() *JsonCodec {
	codec := &JsonCodec{}
	return codec
}

type JsonCodec struct {
}

func (self *JsonCodec) Encode(buffer *fairy.Buffer, msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	buffer.Append(data)
	return nil
}

func (self *JsonCodec) Decode(buffer *fairy.Buffer, msg interface{}) error {
	data := buffer.ReadToEnd()
	return json.Unmarshal(data, msg)
}
