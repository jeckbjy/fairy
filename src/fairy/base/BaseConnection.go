package base

import (
	"fairy"
	"sync/atomic"
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
	State     int32
	side      int
	data      interface{}
	Host      string
}

func (self *BaseConnection) NewBase(transport fairy.Transport, filters fairy.FilterChain, side bool, kind int) {
	self.transport = transport
	self.FilterChain = filters
	self.ctype = kind
	if side {
		self.side = SIDE_SERVER
	} else {
		self.side = SIDE_CLIENT
	}
	self.State = fairy.ConnStateClosed
	self.connId = fairy.GetConnMgr().NewId(side)
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

func (self *BaseConnection) IsState(state int32) bool {
	return self.State == int32(state)
}

func (self *BaseConnection) SwapState(old int32, new int32) bool {
	return atomic.CompareAndSwapInt32(&self.State, old, new)
}

func (self *BaseConnection) GetState() int32 {
	return self.State
}

func (self *BaseConnection) SetState(state int32) {
	self.State = state
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

func (self *BaseConnection) GetConfig(key *fairy.AttrKey) interface{} {
	return self.transport.GetConfig(key)
}
