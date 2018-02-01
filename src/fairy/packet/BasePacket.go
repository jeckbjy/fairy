package packet

import (
	"reflect"
)

func NewBasePacket() *BasePacket {
	packet := &BasePacket{}
	return packet
}

type BasePacket struct {
	id       uint
	name     string
	msg      interface{}
	result   uint
	serialId uint
	time     uint
	checksum uint64
}

func (self *BasePacket) GetId() uint {
	return self.id
}

func (self *BasePacket) SetId(id uint) {
	self.id = id
}

func (self *BasePacket) GetName() string {
	return self.name
}

func (self *BasePacket) SetName(name string) {
	self.name = name
}

func (self *BasePacket) GetMessage() interface{} {
	return self.msg
}

func (self *BasePacket) SetMessage(msg interface{}) {
	self.msg = msg
	if msg != nil && self.name == "" {
		self.name = reflect.TypeOf(msg).Name()
	}
}

func (self *BasePacket) GetResult() uint {
	return self.result
}

func (self *BasePacket) SetResult(r uint) {
	self.result = r
}

func (self *BasePacket) GetSerialId() uint {
	return self.serialId
}

func (self *BasePacket) SetSerialId(id uint) {
	self.serialId = id
}

func (self *BasePacket) GetTime() uint {
	return self.time
}

func (self *BasePacket) SetTime(t uint) {
	self.time = t
}

func (self *BasePacket) SetChecksum(v uint64) {
	self.checksum = v
}

func (self *BasePacket) GetChecksum() uint64 {
	return self.checksum
}
