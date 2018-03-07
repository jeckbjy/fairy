package exec

import (
	"fairy"
)

type Event interface {
	Process()
}

type Callback func()

func NewFuncEvent(cb Callback) *FuncEvent {
	ev := &FuncEvent{cb: cb}
	return ev
}

type FuncEvent struct {
	cb Callback
}

func (self *FuncEvent) Process() {
	self.cb()
}

func NewPacketEvent(conn fairy.Conn, packet fairy.Packet, handler fairy.Handler) *PacketEvent {
	ev := &PacketEvent{conn: conn, packet: packet, handler: handler}
	return ev
}

type PacketEvent struct {
	conn    fairy.Conn
	packet  fairy.Packet
	handler fairy.Handler
}

func (self *PacketEvent) Process() {
	self.handler.Invoke(self.conn, self.packet)
}
