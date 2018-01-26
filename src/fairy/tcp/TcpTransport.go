package tcp

import (
	"fairy"
	"fairy/base"
	"net"
	"sync"
)

func NewTransport() fairy.Transport {
	tran := &TcpTransport{}
	tran.NewBase()
	tran.stopFlag = make(chan bool)
	tran.waitGroup = sync.WaitGroup{}
	return tran
}

type TcpTransport struct {
	base.BaseTransport
	stopFlag  chan bool
	waitGroup sync.WaitGroup
}

func (t *TcpTransport) Listen(host string, ctype int) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}

	t.waitGroup.Add(1)
	go func() {
		defer t.waitGroup.Done()

		for {
			select {
			case <-t.stopFlag:
				break
			default:
				conn, err := listener.Accept()
				if err == nil {
					newConn := NewConnection(t, t.GetFilterChain(), true, ctype)
					newConn.Open(conn)
					newConn.HandleOpen(newConn)
				} else {
					fairy.Error("accept fail!")
				}
			}
		}
	}()
}

func (t *TcpTransport) Connect(host string, ctype int) fairy.ConnectFuture {
	newConn := NewConnection(t, t.GetFilterChain(), false, ctype)
	future := base.NewConnectFuture(newConn)
	t.ConnectBy(future, newConn, host)
	return future
}

func (t *TcpTransport) ConnectBy(future fairy.ConnectFuture, newConn *TcpConnection, host string) {
	t.waitGroup.Add(1)
	go func() {
		// wait for close
		defer t.waitGroup.Done()

		conn, err := net.Dial("tcp", host)
		if future == nil || future.Result() != fairy.FUTURE_RESULT_TIMEOUT {
			future_result := 0
			if err == nil {
				newConn.Open(conn)
				newConn.HandleOpen(newConn)
				future_result = fairy.FUTURE_RESULT_SUCCEED
			} else {
				// panic
				newConn.HandleError(newConn, fairy.ErrConnectFail)
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

func (t *TcpTransport) Start() {
	t.waitGroup.Add(1)
}

func (t *TcpTransport) Stop() {
	close(t.stopFlag)
	t.waitGroup.Done()
	t.waitGroup.Wait()
}

func (t *TcpTransport) Wait() {
	t.waitGroup.Wait()
}

func (t *TcpTransport) OnExit() {
	t.Stop()
}

func (t *TcpTransport) Reconnect(conn *TcpConnection) {
	// 断线重连
	if t.CfgReconnectInterval == 0 {
		t.ConnectBy(nil, conn, conn.Host)
	} else {
		fairy.StartTimer(int64(t.CfgReconnectInterval*1000), func(*fairy.Timer) {
			t.ConnectBy(nil, conn, conn.Host)
		})
	}
}

func (t *TcpTransport) HandleConnClose(conn *TcpConnection) {
	if conn.IsClientSide() && t.IsNeedReconnect() {
		t.Reconnect(conn)
	} else {
		// reomve conn??
	}
}
