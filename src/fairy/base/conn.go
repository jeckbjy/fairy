package base

import (
	"fairy"
	"sync/atomic"
)

const (
	SIDE_SERVER = 0
	SIDE_CLIENT = 1
)

type Connection struct {
	AttrMap
	fairy.FilterChain
	transport fairy.Transport
	connId    uint
	openId    string
	uid       uint64
	kind      int
	State     int32
	side      int
	data      interface{}
	Host      string
}

func (self *Connection) NewBase(transport fairy.Transport, filters fairy.FilterChain, side bool, kind int) {
	self.transport = transport
	self.FilterChain = filters
	self.kind = kind
	if side {
		self.side = SIDE_SERVER
	} else {
		self.side = SIDE_CLIENT
	}
	self.State = fairy.ConnStateClosed
	self.connId = fairy.GetConnMgr().NewId(side)
}

func (self *Connection) GetType() int {
	return self.kind
}

func (self *Connection) SetType(ctype int) {
	self.kind = ctype
}

func (self *Connection) GetConnId() uint {
	return self.connId
}

func (self *Connection) SetConnId(id uint) {
	self.connId = id
}

func (self *Connection) GetUid() uint64 {
	return self.uid
}

func (self *Connection) SetUid(uid uint64) {
	self.uid = uid
}

func (self *Connection) GetOpenId() string {
	return self.openId
}

func (self *Connection) SetOpenId(id string) {
	self.openId = id
}

func (self *Connection) IsState(state int32) bool {
	return self.State == int32(state)
}

func (self *Connection) SwapState(old int32, new int32) bool {
	return atomic.CompareAndSwapInt32(&self.State, old, new)
}

func (self *Connection) GetState() int32 {
	return self.State
}

func (self *Connection) SetState(state int32) {
	self.State = state
}

func (self *Connection) GetData() interface{} {
	return self.data
}

func (self *Connection) SetData(data interface{}) {
	self.data = data
}

func (self *Connection) IsServerSide() bool {
	return self.side == SIDE_SERVER
}

func (self *Connection) IsClientSide() bool {
	return self.side == SIDE_CLIENT
}

func (self *Connection) GetTransport() fairy.Transport {
	return self.transport
}

func (self *Connection) GetConfig(key *fairy.AttrKey) interface{} {
	return self.transport.GetConfig(key)
}
