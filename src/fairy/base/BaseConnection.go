package base

import (
	"fairy"
)

const (
	SIDE_SERVER = 0
	SIDE_CLIENT = 1
)

type BaseConnection struct {
	BaseAttrMap
	fairy.FilterChain
	transport fairy.Transport
	connId    uint
	openId    string
	uid       uint64
	ctype     int
	state     int
	side      int
	data      interface{}
}

func (self *BaseConnection) Create(transport fairy.Transport, filters fairy.FilterChain, side bool, kind int) {
	self.transport = transport
	self.FilterChain = filters
	self.ctype = kind
	if side {
		self.side = SIDE_SERVER
	} else {
		self.side = SIDE_CLIENT
	}
}

func (self *BaseConnection) GetType() int {
	return self.ctype
}

func (self *BaseConnection) SetType(ctype int) {
	self.ctype = ctype
}

func (self *BaseConnection) GetConnId() uint {
	return self.connId
}

func (self *BaseConnection) SetConnId(id uint) {
	self.connId = id
}

func (self *BaseConnection) GetUid() uint64 {
	return self.uid
}

func (self *BaseConnection) SetUid(uid uint64) {
	self.uid = uid
}

func (self *BaseConnection) GetOpenId() string {
	return self.openId
}

func (self *BaseConnection) SetOpenId(id string) {
	self.openId = id
}

func (self *BaseConnection) IsState(state int) bool {
	return self.state == state
}

func (self *BaseConnection) GetState() int {
	return self.state
}

func (self *BaseConnection) SetState(state int) {
	self.state = state
}

func (self *BaseConnection) GetData() interface{} {
	return self.data
}

func (self *BaseConnection) SetData(data interface{}) {
	self.data = data
}

func (self *BaseConnection) IsServerSide() bool {
	return self.side == SIDE_SERVER
}

func (self *BaseConnection) IsClientSide() bool {
	return self.side == SIDE_CLIENT
}

func (self *BaseConnection) GetTransport() fairy.Transport {
	return self.transport
}
