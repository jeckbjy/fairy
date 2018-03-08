package kcp

import (
	"fairy"
	"fairy/base"
	"fairy/snet"
	"fairy/timer"
	"fmt"
	"net"

	kcpgo "github.com/xtaci/kcp-go"
)

func NewTransport() fairy.Transport {
	tran := &KcpTran{}
	tran.Create()
	return tran
}

type KcpTran struct {
	snet.StreamTran
}

func (kt *KcpTran) Listen(host string, kind int) error {
	listener, err := kcpgo.ListenWithOptions(host, nil, 10, 3)
	if err != nil {
		return err
	}

	kt.AddListener(listener)

	//
	kt.AddGroup()
	go func(kt *KcpTran, listener net.Listener, kind int) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			kcp_conn := NewConn(kt, true, kind)
			kcp_conn.Open(conn)
		}
		kt.Done()
	}(kt, listener, kind)

	return nil
}

func (kt *KcpTran) ConnectBy(promise fairy.Promise, new_conn fairy.Conn) (fairy.Future, error) {
	kt.AddGroup()
	stream_conn := new_conn.(*snet.StreamConn)
	host := stream_conn.GetHost()

	go func() {
		conn, err := kcpgo.DialWithOptions(host, nil, 10, 3)
		if !promise.IsCanceled() {
			if err == nil {
				stream_conn.Open(conn)
				promise.SetSuccess()
			} else {
				stream_conn.Error(err)
				promise.SetFailure()
			}
		}

		kt.Done()
	}()

	return promise, nil
}

func (kt *KcpTran) Connect(host string, kind int) (fairy.Future, error) {
	new_conn := NewConn(kt, false, kind)
	new_conn.SetHost(host)
	promise := base.NewPromise(new_conn)
	return kt.ConnectBy(promise, new_conn)
}

func (kt *KcpTran) Reconnect(conn fairy.Conn) (fairy.Future, error) {
	if kt.IsStopped() {
		return nil, fmt.Errorf("stopped, cannot reconnect")
	}

	interval := kt.CfgReconnectInterval

	promise := base.NewPromise(conn)

	if interval == 0 {
		kt.ConnectBy(promise, conn)
	} else {
		timer.Start(int64(interval*1000), func(*timer.Timer) {
			kt.ConnectBy(promise, conn)
		})

	}

	return promise, nil
}
