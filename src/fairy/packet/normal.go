package packet

import (
	"fairy"
)

func NewNormal() *BasePacket {
	pkt := &BasePacket{}
	return pkt
}

func EncodeNormal(pkt *BasePacket, buffer *fairy.Buffer) error {
	writer := NewWriter(buffer)
	writer.PutId(pkt.GetId())
	// TODO:others rpc
	writer.Flush()
	return nil
}

func DecodeNormal(pkt *BasePacket, buffer *fairy.Buffer) error {
	reader := NewReader(buffer)
	id, err := reader.GetId()
	if err != nil {
		return err
	}
	pkt.SetId(id)
	// TODO:others rpc
	return nil
}
