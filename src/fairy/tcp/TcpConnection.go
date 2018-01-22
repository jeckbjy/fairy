package tcp

import (
	"container/list"
	"fairy"
	"fairy/base"
	"net"
	"sync"
)

const (
	WRITE_FLAG_FINISH  = 0 // 写完成
	WRITE_FLAG_WRITING = 1 // 写当中
)

func NewConnection(tran fairy.Transport, filters fairy.FilterChain, side bool, ctype int) *TcpConnection {
	conn := &TcpConnection{}
	conn.Create(tran, filters, side, ctype)
	conn.Init()
	return conn
}

type TcpConnection struct {
	base.BaseConnection
	net.Conn
	stopFlag     chan bool
	waitGroup    sync.WaitGroup
	writerFuture *base.BaseFuture // 用于阻塞写数据完成
	writerLock   *sync.Mutex      // 写锁
	writerCond   *sync.Cond
	writer       *list.List
	reader       *fairy.Buffer
}

func (self *TcpConnection) Init() {
	self.stopFlag = make(chan bool)
	self.reader = fairy.NewBuffer()
	self.writerLock = &sync.Mutex{}

	// lazy init
	// self.writer = list.New()
	// self.writerFuture = base.NewFuture()
	// self.writerCond = sync.NewCond(self.writerLock)
}

func (self *TcpConnection) Send(obj interface{}) {
	self.HandleWrite(self, obj)
}

func (self *TcpConnection) Flush() {
	if self.writerFuture != nil {
		// 阻塞到所有数据写完
		self.writerFuture.Wait(-1)
	}
}

func (self *TcpConnection) Write(buffer *fairy.Buffer) {
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

func (self *TcpConnection) Read() *fairy.Buffer {
	return self.reader
}

func (self *TcpConnection) Open(conn net.Conn) {
	self.Conn = conn
	go self.readThread()
	// go self.sendThread()
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

func (self *TcpConnection) readThread() {
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
				self.HandleError(self, err)
			}
		}
	}
}

func (self *TcpConnection) sendThread() {
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
					_, err = self.Conn.Write(data)
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
