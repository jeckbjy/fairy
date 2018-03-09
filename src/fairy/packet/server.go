package packet

import (
	"fairy"
)

func NewServer() *ServerPacket {
	pkt := &ServerPacket{}
	return pkt
}

// ServerPacket 服务器内部编码,编码规则,flag(1byte)+[head....]
type ServerPacket struct {
	BasePacket
	Mode uint
	Uid  uint64
	Host string
}

func EncodeServer(pkt *ServerPacket, buffer *fairy.Buffer) error {
	writer := NewWriter(buffer)
	// normal
	writer.PutId(pkt.GetId())
	writer.PutUint64(pkt.GetRpcId())
	writer.PutUint(pkt.GetResult())
	// other
	writer.PutUint(pkt.Mode)
	writer.PutUint64(pkt.Uid)
	writer.PutStr(pkt.Host)
	writer.Flush()
	return nil
}

func DecodeServer(pkt *ServerPacket, buffer *fairy.Buffer) error {
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

	if pkt.Mode, err = reader.GetUint(); err != nil {
		return err
	}

	if pkt.Uid, err = reader.GetUint64(); err != nil {
		return err
	}

	if pkt.Host, err = reader.GetStr(); err != nil {
		return err
	}

	return nil
}
