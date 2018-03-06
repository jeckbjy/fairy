package tcp

import (
	"fairy"
	"fairy/snet"
	"net"
)

func NewConn(tran fairy.Transport, side bool, kind int) *snet.StreamConn {
	conn := snet.NewConn(&TcpConn{}, tran, side, kind)
	return conn
}

type TcpConn struct {
	net.Conn
}

func (tc *TcpConn) Open(conn interface{}) {
	tc.Conn = conn.(net.Conn)
}

func (tc *TcpConn) Read(cap int) ([]byte, error) {
	data := make([]byte, cap)
	n, err := tc.Conn.Read(data)
	if err != nil {
		return nil, err
	}

	return data[:n], nil
}

func (tc *TcpConn) Write(buf []byte) error {
	_, err := tc.Conn.Write(buf)
	return err
}

// type ConnAdapter struct {
// 	net.Conn
// }

// func (c *ConnAdapter) Read() error {

// 	return nil
// }

// func (c *ConnAdapter) Write() error {
// 	return nil
// }

// func NewConnection(tran fairy.Transport, filters fairy.FilterChain, side bool, ctype int) *TcpConn {
// 	conn := &TcpConn{}
// 	conn.Create(tran, filters, side, ctype)
// 	return conn
// }

// type TcpConn struct {
// 	base.Conn
// 	base.ConnWriter
// 	base.ConnReader
// 	conn net.Conn
// 	wg   sync.WaitGroup
// }

// func (self *TcpConn) Open(conn net.Conn) {
// 	if self.SwapState(fairy.ConnStateClosed, fairy.ConnStateOpen) {
// 		self.conn = conn
// 		self.NewWriter()
// 		self.NewReader()
// 		fairy.GetConnMgr().Put(self)
// 		go self.readThread()
// 		self.HandleOpen(self)
// 	}
// }

// func (self *TcpConn) Close() {
// 	// 线程安全调用
// 	if self.SwapState(fairy.ConnStateOpen, fairy.ConnStateConnecting) {
// 		// 异步关闭，需要等待读写线程退出，才能退出
// 		go func() {
// 			self.HandleClose(self)
// 			self.conn.Close()
// 			self.StopWrite()
// 			self.wg.Wait()
// 			self.SetState(fairy.ConnStateClosed)
// 			self.conn = nil
// 			// remove
// 			fairy.GetConnMgr().Remove(self.GetConnId())
// 			// try reconnect
// 			trans := self.GetTransport().(*TcpTran)
// 			trans.Reconnect(self)
// 		}()
// 	}
// }

// func (self *TcpConn) Read() *fairy.Buffer {
// 	return self.Reader()
// }

// func (self *TcpConn) Write(buffer *fairy.Buffer) {
// 	self.PushBuffer(buffer, self.sendThread)
// }

// func (self *TcpConn) Send(msg interface{}) {
// 	self.HandleWrite(self, msg)
// }

// func (self *TcpConn) readThread() {
// 	self.wg.Add(1)
// 	defer self.wg.Done()
// 	// loop read
// 	bufferSize := self.GetConfig(fairy.KeyReaderBufferSize).(int)
// 	for {
// 		// 读取数据
// 		data := make([]byte, bufferSize)
// 		n, err := self.conn.Read(data)
// 		if err == nil {
// 			self.Append(data[:n])
// 			self.HandleRead(self)
// 		} else {
// 			self.HandleError(self, err)
// 			break
// 		}
// 	}
// }

// func (self *TcpConn) sendThread() {
// 	self.wg.Add(1)
// 	defer self.wg.Done()

// 	for !self.IsStopped() {
// 		buffers := list.List{}
// 		self.WaitBuffers(&buffers)

// 		// write all buffer
// 		err := self.WriteBuffers(self.conn, &buffers)
// 		if err != nil {
// 			self.HandleError(self, err)
// 			break
// 		}
// 	}
// }

// func (self *TcpConn) LocalAddr() net.Addr {
// 	return self.conn.LocalAddr()
// }

// func (self *TcpConn) RemoteAddr() net.Addr {
// 	return self.conn.RemoteAddr()
// }
