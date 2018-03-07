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
	writer.PutUint64(pkt.GetRpcId())
	writer.PutUint(pkt.GetResult())
	writer.Flush()
	return nil
}

func DecodeNormal(pkt *BasePacket, buffer *fairy.Buffer) error {
	reader := NewReader(buffer)
	var err error
	var id uint
	if id, err = reader.GetId(); err != nil {
		return err
	}
	pkt.SetId(id)

	if pkt.rpcid, err = reader.GetUint64(); err != nil {
		return err
	}

	if pkt.result, err = reader.GetUint(); err != nil {
		return err
	}
	return nil
}
