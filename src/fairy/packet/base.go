package packet

import (
	"fairy"
	"reflect"
)

func NewBase() *BasePacket {
	packet := &BasePacket{}
	return packet
}

type BasePacket struct {
	id     uint
	name   string
	msg    interface{}
	result uint
	rpcid  uint64
}

func (self *BasePacket) Reset() {
	self.id = 0
	self.name = ""
	self.msg = nil
	self.result = 0
	self.rpcid = 0
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

func (self *BasePacket) GetRpcId() uint64 {
	return self.rpcid
}

func (self *BasePacket) SetRpcId(id uint64) {
	self.rpcid = id
}

func (self *BasePacket) GetResult() uint {
	return self.result
}

func (self *BasePacket) SetResult(r uint) {
	self.result = r
}

func (self *BasePacket) SetTimeout() {
	self.SetResult(fairy.PacketResultTimeout)
}

func (self *BasePacket) SetSuccess() {
	self.SetResult(fairy.PacketResultSuccess)
}

func (self *BasePacket) SetFailure() {
	self.SetResult(fairy.PacketResultFailure)
}

func (self *BasePacket) IsTimeout() bool {
	return self.result == fairy.PacketResultTimeout
}

func (self *BasePacket) IsSuccess() bool {
	return self.result == fairy.PacketResultSuccess
}

func (self *BasePacket) IsFailure() bool {
	return self.result == fairy.PacketResultFailure
}
