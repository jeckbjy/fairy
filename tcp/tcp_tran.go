package tcp

import (
	"net"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
)

// NewTran create tcp tran
func NewTran() fairy.ITran {
	return base.NewTran(&tcpTran{})
}

type tcpTran struct {
}

func (tt *tcpTran) Connect(host string, options ...fairy.Option) (net.Conn, error) {
	return net.Dial("tcp", host)
}

func (tt *tcpTran) Listen(host string, options ...fairy.Option) (net.Listener, error) {
	return net.Listen("tcp", host)
}

func (tt *tcpTran) Serve(l net.Listener, cb base.OnAccept) {
	for {
		conn, err := l.Accept()
		cb(conn, err)
	}
}

// type TcpTran struct {
// 	base.StreamTran
// }

// func (tt *TcpTran) Listen(host string, tag string) error {
// 	listener, err := net.Listen("tcp", host)
// 	if err != nil {
// 		return err
// 	}

// 	tt.AddListener(listener)
// 	tt.Add()

// 	go func() {
// 		for {
// 			conn, err := listener.Accept()
// 			if err != nil {
// 				break
// 			}

// 			tconn := base.NewConn(tt, false, tag)
// 			tconn.Open(conn)
// 		}

// 		tt.Done()
// 	}()

// 	return nil
// }

// func (tt *TcpTran) Connect(host string, tag string) (fairy.Future, error) {
// 	tconn := base.NewConn(tt, true, tag)
// 	tconn.SetHost(host)
// 	// new_conn := newConn(t, false, tag)
// 	// new_conn.SetHost(host)
// 	// promise := base.NewPromise(new_conn)
// 	// return t.ConnectBy(promise, new_conn)
// }

// func NewTran() fairy.Tran {
// 	tran := &zTcpTran{}
// 	tran.Create()
// 	return tran
// }

// type zTcpTran struct {
// 	snet.StreamTran
// }

// func (t *zTcpTran) Listen(host string, tag interface{}) error {
// 	listener, err := net.Listen("tcp", host)
// 	if err != nil {
// 		log.Error("%+v", err)
// 		return err
// 	}

// 	t.AddListener(listener)

// 	t.AddGroup()
// 	go func(listener net.Listener, tag interface{}) {
// 		for {
// 			conn, err := listener.Accept()
// 			if err != nil {
// 				break
// 			}

// 			new_conn := newConn(t, true, tag)
// 			new_conn.Open(conn)
// 		}

// 		t.Done()
// 	}(listener, tag)

// 	return nil
// }

// func (t *zTcpTran) ConnectBy(promise fairy.Promise, new_conn fairy.Conn) (fairy.Future, error) {
// 	t.AddGroup()
// 	stream_conn := new_conn.(*snet.StreamConn)

// 	go func() {
// 		conn, err := net.Dial("tcp", stream_conn.GetHost())

// 		if !promise.IsCanceled() {
// 			if err == nil {
// 				stream_conn.Open(conn)
// 				promise.SetSuccess()
// 			} else {
// 				stream_conn.Error(err)
// 				promise.SetFailure()
// 			}
// 		}

// 		t.Done()
// 	}()

// 	return promise, nil
// }

// func (t *zTcpTran) Connect(host string, tag interface{}) (fairy.Future, error) {
// 	new_conn := newConn(t, false, tag)
// 	new_conn.SetHost(host)
// 	promise := base.NewPromise(new_conn)
// 	return t.ConnectBy(promise, new_conn)
// }

// func (t *zTcpTran) Reconnect(conn fairy.Conn) (fairy.Future, error) {
// 	if t.IsStopped() {
// 		return nil, fmt.Errorf("stopped, cannot reconnect")
// 	}

// 	interval := t.CfgReconnectInterval

// 	promise := base.NewPromise(conn)

// 	if interval == 0 {
// 		t.ConnectBy(promise, conn)
// 	} else {
// 		timer.Start(timer.ModeDelay, int64(interval*1000), func() {
// 			t.ConnectBy(promise, conn)
// 		})

// 	}
// 	return promise, nil
// }
