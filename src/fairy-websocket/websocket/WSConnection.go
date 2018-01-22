package websocket

import (
	"container/list"
	"fairy"
	"fairy/base"
	"sync"

	"github.com/gorilla/websocket"
)

func NewConnection(tran fairy.Transport, filters fairy.FilterChain, serverSide bool, ctype int) *WSConnection {
	conn := &WSConnection{}
	conn.Create(tran, filters, serverSide, ctype)
	conn.Init()
	return conn
}

type WSConnection struct {
	base.BaseConnection
	*websocket.Conn
	stopFlag     chan bool
	waitGroup    sync.WaitGroup
	reader       *fairy.Buffer
	writer       *list.List
	writerLock   *sync.Mutex // 写锁
	writerCond   *sync.Cond
	writerFuture *base.BaseFuture // 用于阻塞写数据完成
}

func (self *WSConnection) Init() {
	self.reader = fairy.NewBuffer()
	self.writerLock = &sync.Mutex{}
	self.stopFlag = make(chan bool)

	// lazy init when write
	self.writer = nil
	self.writerFuture = nil
	self.writerCond = nil
}

func (self *WSConnection) Open(conn *websocket.Conn) {
	self.Conn = conn
	go self.readThread()
	// go self.sendThread()
}

func (self *WSConnection) Read() *fairy.Buffer {
	return self.reader
}

func (self *WSConnection) Write(buffer *fairy.Buffer) {
	self.writerLock.Lock()
	// check delay run write loop
	if self.writer == nil {
		self.writer = list.New()
		self.writerCond = sync.NewCond(self.writerLock)
		self.writerFuture = base.NewFuture()
		go self.sendThread()
	}
	self.writer.PushBack(buffer)
	self.writerFuture.Reset()

	self.writerCond.Signal()
	self.writerLock.Unlock()
}

func (self *WSConnection) Flush() {
	if self.writerFuture != nil {
		// 阻塞到所有数据写完
		self.writerFuture.Wait(-1)
	}
}

func (self *WSConnection) Send(obj interface{}) {
	self.HandleWrite(self, obj)
}

func (self *WSConnection) Close() fairy.Future {
	future := base.NewFuture()
	self.SetState(fairy.CONN_STATE_CLOSING)
	if self.Conn != nil {
		go func(conn *websocket.Conn) {
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

func (self *WSConnection) readThread() {
	for {
		select {
		case <-self.stopFlag:
			break
		default:
			mtype, data, err := self.Conn.ReadMessage()
			if err != nil {
				switch mtype {
				case websocket.TextMessage, websocket.BinaryMessage:
					self.reader.Append(data)
					self.HandleRead(self)
				case websocket.CloseMessage:
					self.HandleClose(self)
				}
			} else {
				self.HandleError(self, err)
			}
		}
	}
}

func (self *WSConnection) sendThread() {
	// 发送
	self.waitGroup.Add(1)
	defer self.waitGroup.Done()
	// loop write
	for {
		select {
		case <-self.stopFlag:
			break
		default:
			bufferList := list.List{}
			// write buffer,轮询
			self.writerLock.Lock()
			for self.writer.Len() == 0 {
				// 会先unlock,再lock
				self.writerFuture.DoneSucceed()
				self.writerCond.Wait()
			}

			// swap buffer
			bufferList = *self.writer
			self.writer.Init()
			self.writerLock.Unlock()
			// write all buffer
			var err error
			for iter := bufferList.Front(); iter != nil; iter = iter.Next() {
				buffer := iter.Value.(*fairy.Buffer)
				for buffIter := buffer.Front(); buffIter != nil; buffIter = buffIter.Next() {
					data := buffIter.Value.([]byte)
					err = self.Conn.WriteMessage(websocket.BinaryMessage, data)
					if err != nil {
						self.HandleError(self, err)
						break
					}
				}

				if err != nil {
					break
				}
			}
		}
	}
}
