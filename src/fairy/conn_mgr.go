package fairy

import (
	"fairy/util"
	"sync"
)

const maxConnId = 1 << 31

type ConnectionMgr struct {
	conns    map[uint]Conn
	serverId uint
	clientId uint
	mu       sync.Mutex
}

func (self *ConnectionMgr) NewId(serverSide bool) uint {
	self.mu.Lock()
	id := self.genId(serverSide)
	self.mu.Unlock()
	return id
}

func (self *ConnectionMgr) genId(serverSide bool) uint {
	var id uint
	if serverSide {
		self.serverId++
		if self.serverId == maxConnId {
			self.serverId = 1
		}
		id = self.serverId << 1
	} else {
		self.clientId++
		if self.clientId == maxConnId {
			self.clientId = 1
		}
		id = (self.clientId << 1) + 1
	}
	return id
}

func (self *ConnectionMgr) Put(conn Conn) {
	self.mu.Lock()
	if conn.GetConnId() == 0 {
		conn.SetConnId(self.genId(conn.IsServerSide()))
	}

	self.conns[conn.GetConnId()] = conn
	self.mu.Unlock()
}

func (self *ConnectionMgr) Get(id uint) Conn {
	self.mu.Lock()
	defer self.mu.Unlock()
	return self.conns[id]
}

func (self *ConnectionMgr) Remove(id uint) {
	self.mu.Lock()
	delete(self.conns, id)
	self.mu.Unlock()
}

func (self *ConnectionMgr) Close() {
	self.mu.Lock()
	for _, conn := range self.conns {
		conn.Close()
	}
	self.conns = make(map[uint]Conn)
	self.mu.Unlock()
}

var gConnMgr *ConnectionMgr = nil

func GetConnMgr() *ConnectionMgr {
	util.Once(gConnMgr, func() {
		gConnMgr = NewConnMgr()
	})

	return gConnMgr
}

func NewConnMgr() *ConnectionMgr {
	mgr := &ConnectionMgr{}
	mgr.conns = make(map[uint]Conn)
	return mgr
}
