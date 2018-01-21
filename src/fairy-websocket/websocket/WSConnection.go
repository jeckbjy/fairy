package websocket

import (
	"fairy"
	"fairy/base"
	"net"
)

type WSConnection struct {
	base.BaseConnection
}

func (self *WSConnection) LocalAddr() net.Addr {
	return nil
}

func (self *WSConnection) RemoveAddr() net.Addr {
	return nil
}

func (self *WSConnection) Read() *fairy.Buffer {
	return nil
}

func (self *WSConnection) Write(buffer *fairy.Buffer) {

}

func (self *WSConnection) Close() {

}
