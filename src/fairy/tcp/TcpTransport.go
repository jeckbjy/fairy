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
	tran.wg = sync.WaitGroup{}
	return tran
}

type TcpTransport struct {
	base.BaseTransport
	wg        sync.WaitGroup
	listeners []net.Listener
}

func (t *TcpTransport) Listen(host string, ctype int) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	t.wg.Add(1)

	t.listeners = append(t.listeners, listener)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}

			newConn := NewConnection(t, t.GetFilterChain(), true, ctype)
			newConn.Open(conn)
			newConn.HandleOpen(newConn)
		}

		t.wg.Done()
	}()

	return nil
}

func (t *TcpTransport) Connect(host string, ctype int) (fairy.ConnectFuture, error) {
	newConn := NewConnection(t, t.GetFilterChain(), false, ctype)
	newConn.Host = host
	future := base.NewConnectFuture(newConn)
	t.ConnectBy(future, newConn, host)
	return future, nil
}

func (t *TcpTransport) ConnectBy(future fairy.ConnectFuture, newConn *TcpConnection, host string) {
	t.wg.Add(1)
	go func() {
		// wait for close
		defer t.wg.Done()

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

func (t *TcpTransport) TryReconnect(conn *TcpConnection) bool {
	if conn.IsClientSide() && t.IsNeedReconnect() {
		// 断线重连
		if t.CfgReconnectInterval == 0 {
			t.ConnectBy(nil, conn, conn.Host)
		} else {
			fairy.StartTimer(int64(t.CfgReconnectInterval*1000), func(*fairy.Timer) {
				t.ConnectBy(nil, conn, conn.Host)
			})
		}
		return true
	}
	return false
}

func (t *TcpTransport) Start() {
}

func (t *TcpTransport) Stop() {
	// close all listener
	for _, listener := range t.listeners {
		listener.Close()
	}

	t.listeners = nil

	// stop reconnect
	t.CfgReconnectInterval = -1
}

func (t *TcpTransport) Wait() {
	t.wg.Wait()
}

func (t *TcpTransport) OnExit() {
	t.Stop()
}
