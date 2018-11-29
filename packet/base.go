package packet

func NewBase() *BasePacket {
	pkt := &BasePacket{}
	return pkt
}

// BasePacket 最基础的消息包结构
type BasePacket struct {
	id   uint
	name string
	msg  interface{}
}

func (pkt *BasePacket) GetId() uint {
	return pkt.id
}

func (pkt *BasePacket) SetId(id uint) {
	pkt.id = id
}

func (pkt *BasePacket) GetName() string {
	return pkt.name
}

func (pkt *BasePacket) SetName(name string) {
	pkt.name = name
}

func (pkt *BasePacket) GetMessage() interface{} {
	return pkt.msg
}

func (pkt *BasePacket) SetMessage(msg interface{}) {
	pkt.msg = msg
}
