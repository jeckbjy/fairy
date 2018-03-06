package tcp

import (
	"fairy"
	"fairy/base"
	"fairy/log"
	"fairy/snet"
	"fairy/timer"
	"fmt"
	"net"
)

func NewTransport() fairy.Transport {
	tran := &TcpTran{}
	tran.Create()
	return tran
}

type TcpTran struct {
	snet.StreamTran
}

func (t *TcpTran) Listen(host string, kind int) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Error("%+v", err)
		return err
	}

	t.AddListener(listener)

	t.AddGroup()
	go func(listener net.Listener, kind int) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}

			new_conn := NewConn(t, true, kind)
			new_conn.Open(conn)
		}

		t.Done()
	}(listener, kind)

	return nil
}

func (t *TcpTran) ConnectBy(promise fairy.Promise, new_conn fairy.Connection) (fairy.Future, error) {
	t.AddGroup()
	stream_conn := new_conn.(*snet.StreamConn)
	count_max := t.CfgReconnectCount

	go func() {
		ok := false
		for i := 0; i < count_max && !t.IsStopped(); i++ {
			// logging ??
			conn, err := net.Dial("tcp", stream_conn.GetHost())

			if promise.IsCanceled() {
				break
			}

			if err == nil {
				ok = true
				promise.SetSuccess()
				stream_conn.Open(conn)
				break
			} else {
				stream_conn.Error(err)
			}
		}

		if !ok {
			promise.SetFailure()
		}

		t.Done()
	}()

	return promise, nil
}

func (t *TcpTran) Connect(host string, kind int) (fairy.Future, error) {
	new_conn := NewConn(t, false, kind)
	new_conn.SetHost(host)
	promise := base.NewPromise(new_conn)
	return t.ConnectBy(promise, new_conn)
}

func (t *TcpTran) Reconnect(conn fairy.Connection) (fairy.Future, error) {
	if t.IsStopped() {
		return nil, fmt.Errorf("stopped, cannot reconnect")
	}

	interval := t.CfgReconnectInterval

	promise := base.NewPromise(conn)

	if interval == 0 {
		t.ConnectBy(promise, conn)
	} else {
		timer.Start(int64(interval*1000), func(*timer.Timer) {
			t.ConnectBy(promise, conn)
		})

	}
	return promise, nil
}

// func NewTransport() fairy.Transport {
// 	tran := &TcpTran{}
// 	tran.NewBase()
// 	tran.wg = sync.WaitGroup{}
// 	return tran
// }

// type TcpTran struct {
// 	base.Transport
// 	listeners []net.Listener
// 	wg        sync.WaitGroup
// }

// func (t *TcpTran) Listen(host string, kind int) error {
// 	listener, err := net.Listen("tcp", host)
// 	if err != nil {
// 		fairy.Error("%+v", err)
// 		return err
// 	}

// 	t.wg.Add(1)

// 	t.listeners = append(t.listeners, listener)

// 	go func(listener net.Listener, kind int) {
// 		for {
// 			conn, err := listener.Accept()
// 			if err != nil {
// 				break
// 			}

// 			newConn := NewConnection(t, t.GetFilterChain(), true, kind)
// 			newConn.Open(conn)
// 		}

// 		t.wg.Done()
// 	}(listener, kind)

// 	return nil
// }

// func (t *TcpTran) Connect(host string, ctype int) (fairy.ConnectFuture, error) {
// 	newConn := NewConnection(t, t.GetFilterChain(), false, ctype)
// 	newConn.Host = host
// 	future := base.NewConnectFuture(newConn)
// 	t.ConnectBy(future, newConn, host)
// 	return future, nil
// }

// func (t *TcpTran) ConnectBy(future fairy.ConnectFuture, newConn *TcpConn, host string) {
// 	t.wg.Add(1)
// 	go func() {
// 		// wait for close
// 		defer t.wg.Done()

// 		conn, err := net.Dial("tcp", host)
// 		if future == nil || future.Result() != fairy.FUTURE_RESULT_TIMEOUT {
// 			future_result := 0
// 			if err == nil {
// 				newConn.Open(conn)
// 				future_result = fairy.FUTURE_RESULT_SUCCEED
// 			} else {
// 				// panic
// 				newConn.HandleError(newConn, fairy.ErrConnectFail)
// 				future_result = fairy.FUTURE_RESULT_FAIL
// 			}

// 			if future != nil {
// 				future.Done(future_result)
// 			}
// 		} else if err != nil {
// 			conn.Close()
// 		}
// 	}()
// }

// func (t *TcpTran) Reconnect(conn *TcpConn) bool {
// 	if conn.IsClientSide() && t.IsNeedReconnect() {
// 		// 断线重连
// 		if t.CfgReconnectInterval == 0 {
// 			t.ConnectBy(nil, conn, conn.Host)
// 		} else {
// 			fairy.StartTimer(int64(t.CfgReconnectInterval*1000), func(*fairy.Timer) {
// 				t.ConnectBy(nil, conn, conn.Host)
// 			})
// 		}
// 		return true
// 	}
// 	return false
// }

// func (t *TcpTran) Stop() {
// 	// close all listener
// 	for _, listener := range t.listeners {
// 		listener.Close()
// 	}

// 	t.listeners = nil

// 	// stop reconnect
// 	t.CfgReconnectInterval = -1
// }

// func (t *TcpTran) Wait() {
// 	t.wg.Wait()
// }

// func (t *TcpTran) OnExit() {
// 	t.Stop()
// }
