package udp

import (
	"fairy"
	"fairy/base"
	"net"
)

// TODO:
type udpConn struct {
	base.Conn
	*net.UDPConn
	rbuf *fairy.Buffer
}

func (uc *udpConn) Close() {
}

func (uc *udpConn) Read() *fairy.Buffer {
	return uc.rbuf
}

func (uc *udpConn) Write(buf *fairy.Buffer) {

}

func (uc *udpConn) Send(msg interface{}) {
	uc.HandleWrite(uc, msg)
}

func (uc *udpConn) Wait() {
}
