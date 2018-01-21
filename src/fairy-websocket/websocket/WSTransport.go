package websocket

import (
	"fairy"
	"fairy/base"
	"sync"
)

func NewTransport() fairy.Transport {
	ws := &WSTransport{}
	return ws
}

type WSTransport struct {
	base.BaseTransport
	stopFlag  chan bool
	waitGroup sync.WaitGroup
}

func (self *WSTransport) Listen(host string, kind int) {

}

func (self *WSTransport) Connect(host string, kind int) fairy.ConnectFuture {
	return nil
}

func (self *WSTransport) Start() {

}

func (self *WSTransport) Stop() {

}

func (self *WSTransport) Wait() {
	self.waitGroup.Wait()
}

func (self *WSTransport) OnExit() {
	self.Stop()
}
