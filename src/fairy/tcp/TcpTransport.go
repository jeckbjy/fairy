package tcp

import (
	"fairy"
	"fairy/base"
	"net"
	"sync"
)

func NewTransport() fairy.Transport {
	tran := &TcpTransport{}
	tran.BaseTransport.New()
	tran.stopFlag = make(chan bool)
	tran.waitGroup = sync.WaitGroup{}
	return tran
}

type TcpTransport struct {
	base.BaseTransport
	stopFlag  chan bool
	waitGroup sync.WaitGroup
}

func (self *TcpTransport) Listen(host string, ctype int) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
		return
	}

	self.waitGroup.Add(1)
	go func() {
		defer self.waitGroup.Done()

		for {
			select {
			case <-self.stopFlag:
				break
			default:
				conn, err := listener.Accept()
				if err == nil {
					new_conn := NewConnection(self, self.GetFilterChain(), true, ctype)
					new_conn.Open(conn)
					new_conn.HandleOpen(new_conn)
				} else {
					// onError??
				}
			}
		}
	}()
}

func (self *TcpTransport) Connect(host string, ctype int) fairy.ConnectFuture {
	new_conn := NewConnection(self, self.GetFilterChain(), false, ctype)
	future := base.NewConnectFuture(new_conn)
	self.ConnectBy(future, new_conn, host)
	return future
}

func (self *TcpTransport) ConnectBy(future *base.BaseConnectFuture, new_conn *TcpConnection, host string) {
	self.waitGroup.Add(1)
	go func() {
		// wait for close
		defer self.waitGroup.Done()

		conn, err := net.Dial("tcp", host)
		if future == nil || !future.IsResult(fairy.FUTURE_RESULT_TIMEOUT) {
			future_result := 0
			if err != nil {
				new_conn.Open(conn)
				new_conn.HandleOpen(new_conn)
				future_result = fairy.FUTURE_RESULT_SUCCEED
			} else {
				// panic
				new_conn.HandleError(new_conn, fairy.ErrConnectFail)
				future_result = fairy.FUTURE_RESULT_FAIL
			}

			if future != nil {
				future.Done(future_result)
			}
		} else if err != nil {
			conn.Close()
		}
	}()
}

func (self *TcpTransport) Start(waiting bool) {
	self.waitGroup.Add(1)
	if waiting {
		self.waitGroup.Wait()
	}
}

func (self *TcpTransport) Stop() {
	close(self.stopFlag)
	self.waitGroup.Done()
	self.waitGroup.Wait()
}
