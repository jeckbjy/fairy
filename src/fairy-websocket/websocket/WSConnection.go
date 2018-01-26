package websocket

import (
	"container/list"
	"fairy"
	"fairy/base"
	"net"
	"sync"

	"github.com/gorilla/websocket"
)

func NewConnection(tran fairy.Transport, filters fairy.FilterChain, serverSide bool, ctype int) *WSConnection {
	conn := &WSConnection{}
	conn.NewBase(tran, filters, serverSide, ctype)
	return conn
}

type ConnWrapper struct {
	*websocket.Conn
}

func (self *ConnWrapper) Write(data []byte) (int, error) {
	return len(data), self.WriteMessage(websocket.BinaryMessage, data)
}

type WSConnection struct {
	base.BaseConnection
	base.BaseWriter
	base.BaseReader
	conn *websocket.Conn
	wg   sync.WaitGroup
}

func (self *WSConnection) Open(conn *websocket.Conn) {
	if self.SwapState(fairy.ConnStateClosed, fairy.ConnStateOpen) {
		self.conn = conn
		self.NewWriter()
		self.NewReader()
		go self.readThread()
	}
}

func (self *WSConnection) LocalAddr() net.Addr {
	return self.conn.LocalAddr()
}

func (self *WSConnection) RemoteAddr() net.Addr {
	return self.conn.RemoteAddr()
}

func (self *WSConnection) Close() {
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
			trans := self.GetTransport().(*WSTransport)
			trans.HandleConnClose(self)
		}()
	}
}

func (self *WSConnection) Read() *fairy.Buffer {
	return self.Reader()
}

func (self *WSConnection) Write(buffer *fairy.Buffer) {
	self.PushBuffer(buffer, self.sendThread)
}

func (self *WSConnection) Send(msg interface{}) {
	self.HandleWrite(self, msg)
}

func (self *WSConnection) readThread() {
	self.wg.Add(1)
	defer self.wg.Done()
	// loop read
	bufferSize := self.GetConfig(fairy.KeyReaderBufferSize).(int)
	for {
		mtype, data, err := self.conn.ReadMessage()
		if err != nil {
			switch mtype {
			case websocket.TextMessage, websocket.BinaryMessage:
				self.Append(data)
				self.HandleRead(self)
			case websocket.CloseMessage:
				self.Close()
			}
		} else {
			self.HandleError(self, err)
			break
		}
	}
}

func (self *WSConnection) sendThread() {
	self.wg.Add(1)
	defer self.wg.Done()

	wrapper := ConnWrapper{Conn: self.conn}
	for !self.IsStopped() {
		buffers := list.List{}
		self.WaitBuffers(&buffers)

		// write all buffer
		err := self.WriteBuffers(&wrapper, &buffers)
		if err != nil {
			self.HandleError(self, err)
			break
		}
	}
}

// type WSConnection struct {
// 	base.BaseConnection
// 	*websocket.Conn
// 	stopFlag     chan bool
// 	waitGroup    sync.WaitGroup
// 	reader       *fairy.Buffer
// 	writer       *list.List
// 	writerLock   *sync.Mutex // 写锁
// 	writerCond   *sync.Cond
// 	writerFuture *base.BaseFuture // 用于阻塞写数据完成
// }

// func (self *WSConnection) Init() {
// 	self.reader = fairy.NewBuffer()
// 	self.writerLock = &sync.Mutex{}

// 	// lazy init when write
// 	self.stopFlag = nil
// 	self.writer = nil
// 	self.writerFuture = nil
// 	self.writerCond = nil
// }

// func (self *WSConnection) Open(conn *websocket.Conn) {
// 	if conn == nil || self.Conn != nil {
// 		return
// 	}
// 	self.SetState(fairy.ConnStateOpen)
// 	self.Conn = conn
// 	self.stopFlag = make(chan bool)

// 	go self.readThread()
// 	// go self.sendThread()
// }

// func (self *WSConnection) Close() {
// 	if self.Conn != nil {
// 		self.Conn = nil
// 		self.SetState(fairy.ConnStateClosed)
// 		self.Conn.Close()
// 		close(self.stopFlag)
// 		// notify transport
// 		// trans := self.GetTransport().(*TcpTransport)
// 	}
// }

// func (self *WSConnection) Read() *fairy.Buffer {
// 	return self.reader
// }

// func (self *WSConnection) Write(buffer *fairy.Buffer) {
// 	self.writerLock.Lock()
// 	// check delay run write loop
// 	if self.writer == nil {
// 		self.writer = list.New()
// 		self.writerCond = sync.NewCond(self.writerLock)
// 		self.writerFuture = base.NewFuture()
// 		go self.sendThread()
// 	}
// 	self.writer.PushBack(buffer)
// 	self.writerFuture.Reset()

// 	self.writerCond.Signal()
// 	self.writerLock.Unlock()
// }

// func (self *WSConnection) Flush() {
// 	if self.writerFuture != nil {
// 		// 阻塞到所有数据写完
// 		self.writerFuture.Wait(-1)
// 	}
// }

// func (self *WSConnection) Send(obj interface{}) {
// 	self.HandleWrite(self, obj)
// }

// func (self *WSConnection) readThread() {
// 	for {
// 		select {
// 		case <-self.stopFlag:
// 			break
// 		default:
// 			mtype, data, err := self.Conn.ReadMessage()
// 			if err != nil {
// 				switch mtype {
// 				case websocket.TextMessage, websocket.BinaryMessage:
// 					self.reader.Append(data)
// 					self.HandleRead(self)
// 				case websocket.CloseMessage:
// 					self.HandleClose(self)
// 				}
// 			} else {
// 				self.HandleError(self, err)
// 			}
// 		}
// 	}
// }

// func (self *WSConnection) sendThread() {
// 	// 发送
// 	self.waitGroup.Add(1)
// 	defer self.waitGroup.Done()
// 	// loop write
// 	for {
// 		select {
// 		case <-self.stopFlag:
// 			break
// 		default:
// 			bufferList := list.List{}
// 			// write buffer,轮询
// 			self.writerLock.Lock()
// 			for self.writer.Len() == 0 {
// 				// 会先unlock,再lock
// 				self.writerFuture.DoneSucceed()
// 				self.writerCond.Wait()
// 			}

// 			// swap buffer
// 			bufferList = *self.writer
// 			self.writer.Init()
// 			self.writerLock.Unlock()
// 			// write all buffer
// 			var err error
// 			for iter := bufferList.Front(); iter != nil; iter = iter.Next() {
// 				buffer := iter.Value.(*fairy.Buffer)
// 				for buffIter := buffer.Front(); buffIter != nil; buffIter = buffIter.Next() {
// 					data := buffIter.Value.([]byte)
// 					err = self.Conn.WriteMessage(websocket.BinaryMessage, data)
// 					if err != nil {
// 						self.HandleError(self, err)
// 						break
// 					}
// 				}

// 				if err != nil {
// 					break
// 				}
// 			}
// 		}
// 	}
// }
