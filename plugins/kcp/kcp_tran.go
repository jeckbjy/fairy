package kcp

import (
	"fmt"
	"net"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
	"github.com/jeckbjy/fairy/snet"
	"github.com/jeckbjy/fairy/timer"

	kcpgo "github.com/xtaci/kcp-go"
)

func NewTran() fairy.Tran {
	tran := &KcpTran{}
	tran.Create()
	return tran
}

type KcpTran struct {
	snet.StreamTran
}

func (kt *KcpTran) Listen(host string, tag interface{}) error {
	listener, err := kcpgo.ListenWithOptions(host, nil, 10, 3)
	if err != nil {
		return err
	}

	kt.AddListener(listener)

	//
	kt.AddGroup()
	go func(kt *KcpTran, listener net.Listener, tag interface{}) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}
			kcp_conn := NewConn(kt, true, tag)
			kcp_conn.Open(conn)
		}
		kt.Done()
	}(kt, listener, tag)

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

func (kt *KcpTran) Connect(host string, tag interface{}) (fairy.Future, error) {
	new_conn := NewConn(kt, false, tag)
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
		timer.Start(timer.ModeDelay, int64(interval*1000), func() {
			kt.ConnectBy(promise, conn)
		})

	}

	return promise, nil
}
