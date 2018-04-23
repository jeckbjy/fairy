package tcp

import (
	"fmt"
	"net"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
	"github.com/jeckbjy/fairy/log"
	"github.com/jeckbjy/fairy/snet"
	"github.com/jeckbjy/fairy/timer"
)

func NewTransport() fairy.Tran {
	tran := &zTcpTran{}
	tran.Create()
	return tran
}

type zTcpTran struct {
	snet.StreamTran
}

func (t *zTcpTran) Listen(host string, tag interface{}) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Error("%+v", err)
		return err
	}

	t.AddListener(listener)

	t.AddGroup()
	go func(listener net.Listener, tag interface{}) {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}

			new_conn := newConn(t, true, tag)
			new_conn.Open(conn)
		}

		t.Done()
	}(listener, tag)

	return nil
}

func (t *zTcpTran) ConnectBy(promise fairy.Promise, new_conn fairy.Conn) (fairy.Future, error) {
	t.AddGroup()
	stream_conn := new_conn.(*snet.StreamConn)

	go func() {
		conn, err := net.Dial("tcp", stream_conn.GetHost())

		if !promise.IsCanceled() {
			if err == nil {
				stream_conn.Open(conn)
				promise.SetSuccess()
			} else {
				stream_conn.Error(err)
				promise.SetFailure()
			}
		}

		t.Done()
	}()

	return promise, nil
}

func (t *zTcpTran) Connect(host string, tag interface{}) (fairy.Future, error) {
	new_conn := newConn(t, false, tag)
	new_conn.SetHost(host)
	promise := base.NewPromise(new_conn)
	return t.ConnectBy(promise, new_conn)
}

func (t *zTcpTran) Reconnect(conn fairy.Conn) (fairy.Future, error) {
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
