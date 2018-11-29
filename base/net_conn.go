package base

import (
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/sync"
)

const connIdMax = 1 << 31

var connIdLock sync.SpinLock
var connIdCurr uint32

// Conn base connection
type Conn struct {
	AttrMap
	tran      fairy.ITran // transport
	connId    uint        // connection id, auto gen
	tag       string      // user tag
	data      interface{} // user data
	host      string      // used for reconnect
	connector bool        // check is connector or listener
}

func (conn *Conn) Init(tran fairy.ITran, connector bool, tag string) {
	connIdLock.Lock()
	connIdCurr++
	if connIdCurr == connIdMax {
		connIdCurr = 1
	}
	id := connIdCurr << 1
	if connector {
		id++
	}
	conn.connId = uint(id)
	connIdLock.Unlock()

	conn.tran = tran
	conn.connector = connector
	conn.tag = tag
}

func (conn *Conn) GetTran() fairy.ITran {
	return conn.tran
}

func (conn *Conn) GetChain() fairy.IFilterChain {
	return conn.tran.GetChain()
}

func (conn *Conn) GetId() uint {
	return conn.connId
}

func (conn *Conn) GetTag() string {
	return conn.tag
}

func (conn *Conn) SetTag(tag string) {
	conn.tag = tag
}

func (conn *Conn) GetData() interface{} {
	return conn.data
}

func (conn *Conn) SetData(data interface{}) {
	conn.data = data
}

func (conn *Conn) GetHost() string {
	return conn.host
}

func (conn *Conn) SetHost(host string) {
	conn.host = host
}

func (conn *Conn) IsConnector() bool {
	return conn.connector
}

// const (
// 	SIDE_SERVER = 0
// 	SIDE_CLIENT = 1
// )

// type Conn struct {
// 	AttrMap
// 	fairy.FilterChain
// 	tran   fairy.Tran
// 	connId uint
// 	openId string
// 	uid    uint64
// 	kind   int
// 	State  int32
// 	side   int
// 	data   interface{}
// 	tag    interface{}
// 	host   string
// }

// func (self *Conn) Create(tran fairy.Tran, side bool, tag interface{}) {
// 	self.tran = tran
// 	self.FilterChain = tran.GetFilterChain()
// 	self.tag = tag
// 	if side {
// 		self.side = SIDE_SERVER
// 	} else {
// 		self.side = SIDE_CLIENT
// 	}
// 	self.State = fairy.ConnStateClosed
// 	self.connId = fairy.GetConnMgr().NewId(side)
// }

// func (self *Conn) GetType() int {
// 	return self.kind
// }

// func (self *Conn) SetType(ctype int) {
// 	self.kind = ctype
// }

// func (self *Conn) GetConnId() uint {
// 	return self.connId
// }

// func (self *Conn) SetConnId(id uint) {
// 	self.connId = id
// }

// func (self *Conn) GetUid() uint64 {
// 	return self.uid
// }

// func (self *Conn) SetUid(uid uint64) {
// 	self.uid = uid
// }

// func (self *Conn) GetOpenId() string {
// 	return self.openId
// }

// func (self *Conn) SetOpenId(id string) {
// 	self.openId = id
// }

// func (self *Conn) IsState(state int32) bool {
// 	return self.State == int32(state)
// }

// func (self *Conn) SwapState(old int32, new int32) bool {
// 	return atomic.CompareAndSwapInt32(&self.State, old, new)
// }

// func (self *Conn) GetState() int32 {
// 	return self.State
// }

// func (self *Conn) SetState(state int32) {
// 	self.State = state
// }

// func (self *Conn) IsActive() bool {
// 	return self.IsState(fairy.ConnStateOpen)
// }

// func (self *Conn) GetTag() interface{} {
// 	return self.tag
// }

// func (self *Conn) SetTag(tag interface{}) {
// 	self.tag = tag
// }

// func (self *Conn) GetData() interface{} {
// 	return self.data
// }

// func (self *Conn) SetData(data interface{}) {
// 	self.data = data
// }

// func (self *Conn) SetHost(host string) {
// 	self.host = host
// }

// func (self *Conn) GetHost() string {
// 	return self.host
// }

// func (self *Conn) IsServerSide() bool {
// 	return self.side == SIDE_SERVER
// }

// func (self *Conn) IsClientSide() bool {
// 	return self.side == SIDE_CLIENT
// }

// func (self *Conn) GetTransport() fairy.Tran {
// 	return self.tran
// }

// func (self *Conn) GetConfig(key *fairy.AttrKey) interface{} {
// 	return self.tran.GetConfig(key)
// }
