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

			new_conn := newConn(t, true, kind)
			new_conn.Open(conn)
		}

		t.Done()
	}(listener, kind)

	return nil
}

func (t *TcpTran) ConnectBy(promise fairy.Promise, new_conn fairy.Conn) (fairy.Future, error) {
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
	new_conn := newConn(t, false, kind)
	new_conn.SetHost(host)
	promise := base.NewPromise(new_conn)
	return t.ConnectBy(promise, new_conn)
}

func (t *TcpTran) Reconnect(conn fairy.Conn) (fairy.Future, error) {
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
