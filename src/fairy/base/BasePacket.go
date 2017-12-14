package base

import "reflect"

func NewBasePacket() *BasePacket {
	packet := &BasePacket{}
	return packet
}

type BasePacket struct {
	id   uint
	name string
	msg  interface{}
}

func (self *BasePacket) GetId() uint {
	return self.id
}

func (self *BasePacket) GetName() string {
	return self.name
}

func (self *BasePacket) GetMessage() interface{} {
	return self.msg
}

func (self *BasePacket) SetId(id uint) {
	self.id = id
}

func (self *BasePacket) SetName(name string) {
	self.name = name
}

func (self *BasePacket) SetMessage(msg interface{}) {
	self.msg = msg
	if msg != nil && self.name == "" {
		self.name = reflect.TypeOf(msg).Name()
	}
}
