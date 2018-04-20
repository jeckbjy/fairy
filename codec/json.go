package codec

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

func (self *JsonCodec) Encode(msg interface{}, buffer *fairy.Buffer) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	buffer.Append(data)
	return nil
}

func (self *JsonCodec) Decode(msg interface{}, buffer *fairy.Buffer) error {
	data := buffer.ReadToEnd()
	return json.Unmarshal(data, msg)
}
