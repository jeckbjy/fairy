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

func NewConnection(tran fairy.Transport, filters fairy.FilterChain, serverSide bool, ctype int) *TcpConnection {
	tcp_conn := &TcpConnection{}
	tcp_conn.BaseConnection.New(tran, filters, serverSide)
	tcp_conn.SetType(ctype)
	tcp_conn.reader = fairy.NewBuffer()
	tcp_conn.writerLock = &sync.Mutex{}
	return tcp_conn
}

// todo:读，写限流控制??
type TcpConnection struct {
	base.BaseConnection
	net.Conn
	writerFuture *base.BaseFuture // 用于阻塞写数据完成
	writerLock   *sync.Mutex      // 写锁
	writerCond   *sync.Cond
	writer       *list.List
	reader       *fairy.Buffer
	stopFlag     chan bool
	waitGroup    sync.WaitGroup
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
		go self.runWriteLoop()
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
			for iter := bufferList.Front(); iter != nil; iter = iter.Next() {
				buffer := iter.Value.(*fairy.Buffer)
				if err := buffer.SendAll(self.Conn); err != nil {
					self.HandleError(self, err)
					break
				}
			}
		}
	}
}
