package ws

import (
	"fairy"
	"fairy/snet"
	"fmt"

	"github.com/gorilla/websocket"
)

func newConn(tran fairy.Transport, side bool, kind int) *snet.StreamConn {
	conn := snet.NewConn(&WSConn{}, tran, side, kind)
	return conn
}

type WSConn struct {
	*websocket.Conn
}

func (wc *WSConn) Open(conn interface{}) {
	wc.Conn = conn.(*websocket.Conn)
}

func (wc *WSConn) Read(reader *fairy.Buffer, cap int) error {
	mtype, data, err := wc.ReadMessage()
	if err != nil {
		return err
	}

	switch mtype {
	case websocket.TextMessage, websocket.BinaryMessage:
		reader.Append(data)
		return nil
	case websocket.CloseMessage:
		return fmt.Errorf("close")
	}

	return nil
}

func (wc *WSConn) Write(buf []byte) error {
	return wc.WriteMessage(websocket.BinaryMessage, buf)
}

// func NewConnection(tran fairy.Transport, filters fairy.FilterChain, serverSide bool, ctype int) *WSConnection {
// 	conn := &WSConnection{}
// 	conn.Create(tran, filters, serverSide, ctype)
// 	return conn
// }

// type ConnWrapper struct {
// 	*websocket.Conn
// }

// func (self *ConnWrapper) Write(data []byte) (int, error) {
// 	return len(data), self.WriteMessage(websocket.BinaryMessage, data)
// }

// type WSConnection struct {
// 	base.Conn
// 	base.ConnWriter
// 	base.ConnReader
// 	conn *websocket.Conn
// 	wg   sync.WaitGroup
// }

// func (self *WSConnection) LocalAddr() net.Addr {
// 	return self.conn.LocalAddr()
// }

// func (self *WSConnection) RemoteAddr() net.Addr {
// 	return self.conn.RemoteAddr()
// }

// func (self *WSConnection) Open(conn *websocket.Conn) {
// 	if self.SwapState(fairy.ConnStateClosed, fairy.ConnStateOpen) {
// 		self.conn = conn
// 		self.NewWriter()
// 		self.NewReader()
// 		self.HandleOpen(self)
// 		go self.readThread()
// 	}
// }

// func (self *WSConnection) Close() {
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
// 			// try reconnect
// 			trans := self.GetTransport().(*WSTransport)
// 			trans.TryReconnect(self)
// 		}()
// 	}
// }

// func (self *WSConnection) Read() *fairy.Buffer {
// 	return self.Reader()
// }

// func (self *WSConnection) Write(buffer *fairy.Buffer) {
// 	self.PushBuffer(buffer, self.sendThread)
// }

// func (self *WSConnection) Send(msg interface{}) {
// 	self.HandleWrite(self, msg)
// }

// func (self *WSConnection) readThread() {
// 	self.wg.Add(1)
// 	defer self.wg.Done()
// 	// loop read
// 	for {
// 		mtype, data, err := self.conn.ReadMessage()
// 		if err == nil {
// 			switch mtype {
// 			case websocket.TextMessage, websocket.BinaryMessage:
// 				self.Append(data)
// 				self.HandleRead(self)
// 			case websocket.CloseMessage:
// 				self.Close()
// 			}
// 		} else {
// 			self.HandleError(self, err)
// 			break
// 		}
// 	}
// }

// func (self *WSConnection) sendThread() {
// 	self.wg.Add(1)
// 	defer self.wg.Done()

// 	wrapper := ConnWrapper{Conn: self.conn}
// 	for !self.IsStopped() {
// 		buffers := list.List{}
// 		self.WaitBuffers(&buffers)

// 		// write all buffer
// 		err := self.WriteBuffers(&wrapper, &buffers)
// 		if err != nil {
// 			self.HandleError(self, err)
// 			break
// 		}
// 	}
// }
