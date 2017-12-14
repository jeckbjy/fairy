package tcp

import (
	"fairy"
	"fairy/base"
	"net"
	"sync"
)

func NewConnection(tran fairy.Transport, filters fairy.FilterChain, serverSide bool, ctype int) *TcpConnection {
	tcp_conn := &TcpConnection{}
	tcp_conn.BaseConnection.New(tran, filters, serverSide)
	tcp_conn.SetType(ctype)
	return tcp_conn
}

type TcpConnection struct {
	base.BaseConnection
	net.Conn
	reader    *fairy.Buffer
	stopFlag  chan bool
	waitGroup sync.WaitGroup
}

func (self *TcpConnection) Flush() {

}

func (self *TcpConnection) Write(buffer *fairy.Buffer) {

}

func (self *TcpConnection) Read() *fairy.Buffer {
	return self.reader
}

func (self *TcpConnection) Open(conn net.Conn) {
	self.Conn = conn
	go self.runReadLoop()
	go self.runWriteLoop()
}

func (self *TcpConnection) Close() fairy.Future {
	future := base.NewFuture()
	self.SetState(fairy.CONN_STATE_CLOSING)
	if self.Conn != nil {
		go func(conn net.Conn) {
			conn.Close()
			close(self.stopFlag)
			self.waitGroup.Wait()
			future.Done(fairy.FUTURE_RESULT_SUCCEED)
			self.SetState(fairy.CONN_STATE_CLOSED)
		}(self.Conn)
		self.Conn = nil
	}

	return future
}

func (self *TcpConnection) runReadLoop() {
	self.waitGroup.Add(1)
	defer self.waitGroup.Done()
	// loop read
	for {
		select {
		case <-self.stopFlag:
			break
		default:
			// 读取数据
			data := make([]byte, 1024)
			n, err := self.Conn.Read(data)
			if err == nil {
				self.reader.Append(data[:n])
				self.HandleRead(self)
			} else {
				self.HandleError(self, fairy.ErrReadFail)
			}
		}
	}
}

func (self *TcpConnection) runWriteLoop() {
	// 发送
	self.waitGroup.Add(1)
	defer self.waitGroup.Done()
	// loop write
	for {
		select {
		case <-self.stopFlag:
			break
		default:
			// write buffer,轮询
		}
	}
}
