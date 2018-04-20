package fairy

type Event interface {
	Process()
}

//////////////////////////////////////////////////
// FuncEvent
//////////////////////////////////////////////////

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

//////////////////////////////////////////////////
// PacketEvent
//////////////////////////////////////////////////

func NewPacketEvent(conn Conn, packet Packet, handler Handler) *PacketEvent {
	ev := &PacketEvent{conn: conn, packet: packet, handler: handler}
	return ev
}

type PacketEvent struct {
	conn    Conn
	packet  Packet
	handler Handler
}

func (self *PacketEvent) Process() {
	self.handler.Invoke(self.conn, self.packet)
}
