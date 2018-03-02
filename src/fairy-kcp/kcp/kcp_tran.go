package kcp

import (
	"fairy"
	"fairy/base"
	"net"
	"sync"

	kcpgo "github.com/xtaci/kcp-go"
)

func NewTransport() fairy.Transport {
	tran := &KcpTran{}
	tran.NewBase()
	tran.wg = sync.WaitGroup{}
	return tran
}

type KcpTran struct {
	base.Transport
	listeners []net.Listener
	wg        sync.WaitGroup
}

func (kt *KcpTran) Listen(host string, kind int) error {
	listener, err := kcpgo.ListenWithOptions(host, nil, 10, 3)
	if err != nil {
		return err
	}

	kt.listeners = append(kt.listeners, listener)
	kt.wg.Add(1)
	go func(kt *KcpTran, listener net.Listener, kind int) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			kcp_conn := NewConn(kt, kt.GetFilterChain(), true, kind)
			kcp_conn.Open(conn)
		}
		kt.wg.Done()
	}(kt, listener, kind)

	return nil
}

func (kt *KcpTran) Connect(host string, kind int) (fairy.ConnectFuture, error) {
	kcp_conn := NewConn(kt, kt.GetFilterChain(), false, kind)
	kcp_conn.Host = host
	future := base.NewConnectFuture(kcp_conn)
	kt.ConnectBy(future, kcp_conn, host)
	return future, nil
}

func (kt *KcpTran) ConnectBy(future fairy.ConnectFuture, kconn *KcpConn, host string) {
	kt.wg.Add(1)
	go func() {
		// wait for close
		defer kt.wg.Done()

		conn, err := kcpgo.DialWithOptions(host, nil, 10, 3)
		if future == nil || future.Result() != fairy.FUTURE_RESULT_TIMEOUT {
			future_result := 0
			if err == nil {
				kconn.Open(conn)
				future_result = fairy.FUTURE_RESULT_SUCCEED
			} else {
				kconn.HandleError(kconn, fairy.ErrConnectFail)
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

func (kt *KcpTran) Reconnect(conn *KcpConn) bool {
	if conn.IsClientSide() && kt.IsNeedReconnect() {
		// 断线重连
		if kt.CfgReconnectInterval == 0 {
			kt.ConnectBy(nil, conn, conn.Host)
		} else {
			fairy.StartTimer(int64(kt.CfgReconnectInterval*1000), func(*fairy.Timer) {
				kt.ConnectBy(nil, conn, conn.Host)
			})
		}
		return true
	}
	return false
}

func (kt *KcpTran) Stop() {
	// close all listener
	for _, listener := range kt.listeners {
		listener.Close()
	}

	kt.listeners = nil
	kt.CfgReconnectInterval = -1
}

func (kt *KcpTran) Wait() {
	kt.wg.Wait()
}

func (kt *KcpTran) OnExit() {
	kt.Stop()
}
