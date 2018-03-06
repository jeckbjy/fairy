package base

import (
	"fairy"
	"sync/atomic"
)

const (
	SIDE_SERVER = 0
	SIDE_CLIENT = 1
)

type Conn struct {
	AttrMap
	fairy.FilterChain
	tran   fairy.Transport
	connId uint
	openId string
	uid    uint64
	kind   int
	State  int32
	side   int
	data   interface{}
	host   string
}

func (self *Conn) Create(tran fairy.Transport, side bool, kind int) {
	self.tran = tran
	self.FilterChain = tran.GetFilterChain()
	self.kind = kind
	if side {
		self.side = SIDE_SERVER
	} else {
		self.side = SIDE_CLIENT
	}
	self.State = fairy.ConnStateClosed
	self.connId = fairy.GetConnMgr().NewId(side)
}

func (self *Conn) GetType() int {
	return self.kind
}

func (self *Conn) SetType(ctype int) {
	self.kind = ctype
}

func (self *Conn) GetConnId() uint {
	return self.connId
}

func (self *Conn) SetConnId(id uint) {
	self.connId = id
}

func (self *Conn) GetUid() uint64 {
	return self.uid
}

func (self *Conn) SetUid(uid uint64) {
	self.uid = uid
}

func (self *Conn) GetOpenId() string {
	return self.openId
}

func (self *Conn) SetOpenId(id string) {
	self.openId = id
}

func (self *Conn) IsState(state int32) bool {
	return self.State == int32(state)
}

func (self *Conn) SwapState(old int32, new int32) bool {
	return atomic.CompareAndSwapInt32(&self.State, old, new)
}

func (self *Conn) GetState() int32 {
	return self.State
}

func (self *Conn) SetState(state int32) {
	self.State = state
}

func (self *Conn) GetData() interface{} {
	return self.data
}

func (self *Conn) SetData(data interface{}) {
	self.data = data
}

func (self *Conn) SetHost(host string) {
	self.host = host
}

func (self *Conn) GetHost() string {
	return self.host
}

func (self *Conn) IsServerSide() bool {
	return self.side == SIDE_SERVER
}

func (self *Conn) IsClientSide() bool {
	return self.side == SIDE_CLIENT
}

func (self *Conn) GetTransport() fairy.Transport {
	return self.tran
}

func (self *Conn) GetConfig(key *fairy.AttrKey) interface{} {
	return self.tran.GetConfig(key)
}
