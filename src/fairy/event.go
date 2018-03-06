package fairy

import "fairy/log"

type Event interface {
	Process()
}

type Callback func()

////////////////////////////////////////////////////////////
// FuncEvent
////////////////////////////////////////////////////////////
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

////////////////////////////////////////////////////////////
// TimerEvent
////////////////////////////////////////////////////////////
func NewTimerEvent(e *TimerEngine) *TimerEvent {
	ev := &TimerEvent{engine: e}
	return ev
}

type TimerEvent struct {
	engine *TimerEngine
}

func (self *TimerEvent) Process() {
	self.engine.Invoke()
}

////////////////////////////////////////////////////////////
// PacketEvent
////////////////////////////////////////////////////////////
func NewPacketEvent(conn Connection, packet Packet, handler Handler) *PacketEvent {
	ev := &PacketEvent{conn: conn, packet: packet, handler: handler}
	return ev
}

type PacketEvent struct {
	conn    Connection
	packet  Packet
	handler Handler
}

func (self *PacketEvent) Process() {
	defer log.Catch()
	self.handler.Invoke(self.conn, self.packet)
}
