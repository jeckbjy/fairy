package tcp

import (
	"container/list"
	"fairy"
	"fairy/base"
	"net"
	"sync"
)

func NewConnection(tran fairy.Transport, filters fairy.FilterChain, side bool, ctype int) *TcpConnection {
	conn := &TcpConnection{}
	conn.NewBase(tran, filters, side, ctype)
	return conn
}

type TcpConnection struct {
	base.BaseConnection
	base.BaseWriter
	base.BaseReader
	conn net.Conn
	wg   sync.WaitGroup
}

func (self *TcpConnection) Open(conn net.Conn) {
	if self.SwapState(fairy.ConnStateClosed, fairy.ConnStateOpen) {
		self.conn = conn
		self.NewWriter()
		self.NewReader()
		go self.readThread()
	}
}

func (self *TcpConnection) LocalAddr() net.Addr {
	return self.conn.LocalAddr()
}

func (self *TcpConnection) RemoteAddr() net.Addr {
	return self.conn.RemoteAddr()
}

func (self *TcpConnection) Close() {
	// 线程安全调用
	if self.SwapState(fairy.ConnStateOpen, fairy.ConnStateConnecting) {
		// 异步关闭，需要等待读写线程退出，才能退出
		go func() {
			self.HandleClose(self)
			self.conn.Close()
			self.StopWrite()
			self.wg.Wait()
			self.SetState(fairy.ConnStateClosed)
			self.conn = nil
			// try reconnect
			trans := self.GetTransport().(*TcpTransport)
			trans.TryReconnect(self)
		}()
	}
}

func (self *TcpConnection) Read() *fairy.Buffer {
	return self.Reader()
}

func (self *TcpConnection) Write(buffer *fairy.Buffer) {
	self.PushBuffer(buffer, self.sendThread)
}

func (self *TcpConnection) Send(msg interface{}) {
	self.HandleWrite(self, msg)
}

func (self *TcpConnection) readThread() {
	self.wg.Add(1)
	defer self.wg.Done()
	// loop read
	bufferSize := self.GetConfig(fairy.KeyReaderBufferSize).(int)
	for {
		// 读取数据
		data := make([]byte, bufferSize)
		n, err := self.conn.Read(data)
		if err == nil {
			self.Append(data[:n])
			self.HandleRead(self)
		} else {
			self.HandleError(self, err)
			break
		}
	}
}

func (self *TcpConnection) sendThread() {
	self.wg.Add(1)
	defer self.wg.Done()

	for !self.IsStopped() {
		buffers := list.List{}
		self.WaitBuffers(&buffers)

		// write all buffer
		err := self.WriteBuffers(self.conn, &buffers)
		if err != nil {
			self.HandleError(self, err)
			break
		}
	}
}
