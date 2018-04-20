package codec

import (
	"encoding/xml"

	"github.com/jeckbjy/fairy"
)

func NewXml() *XmlCodec {
	codec := &XmlCodec{}
	return codec
}

type XmlCodec struct {
}

func (self *XmlCodec) Encode(msg interface{}, buffer *fairy.Buffer) error {
	data, err := xml.Marshal(msg)
	if err != nil {
		return err
	}

	buffer.Append(data)
	return nil
}

func (self *XmlCodec) Decode(msg interface{}, buffer *fairy.Buffer) error {
	data := buffer.ReadToEnd()
	return xml.Unmarshal(data, msg)
}
