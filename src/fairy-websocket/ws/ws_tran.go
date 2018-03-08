package ws

import (
	"fairy"
	"fairy/base"
	"fairy/snet"
	"fairy/timer"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

func NewTransport() fairy.Transport {
	tran := &WSTran{}
	tran.Create()
	return tran
}

type ServeHttpHandler struct {
	kind  int
	owner *WSTran
}

func (self *ServeHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := self.owner.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// kind
	new_conn := newConn(self.owner, true, self.kind)
	new_conn.Open(conn)
}

type WSTran struct {
	snet.StreamTran
	websocket.Upgrader
}

func (wt *WSTran) Create() {
	wt.StreamTran.Create()
	wt.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func (wt *WSTran) Listen(host string, kind int) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	wt.AddListener(listener)

	wt.AddGroup()
	go func() {
		for {
			svr := http.Server{Handler: &ServeHttpHandler{kind: kind, owner: wt}}
			err := svr.Serve(listener)
			if err != nil {
				break
			}
		}

		wt.Done()
	}()

	return nil
}

func (wt *WSTran) ConnectBy(promise fairy.Promise, new_conn fairy.Conn) (fairy.Future, error) {
	wt.AddGroup()
	stream_conn := new_conn.(*snet.StreamConn)
	host := stream_conn.GetHost()

	go func() {
		conn, _, err := websocket.DefaultDialer.Dial(host, nil)
		if !promise.IsCanceled() {
			if err == nil {
				stream_conn.Open(conn)
				promise.SetSuccess()
			} else {
				stream_conn.Error(err)
				promise.SetFailure()
			}
		}

		wt.Done()
	}()

	return promise, nil
}

func (wt *WSTran) Connect(host string, kind int) (fairy.Future, error) {
	// convert url:localhost:8888->ws://localhost:8888
	pos := strings.Index(host, "//")
	if pos == -1 {
		host = "ws://" + host
	}

	new_conn := newConn(wt, false, kind)
	new_conn.SetHost(host)
	promise := base.NewPromise(new_conn)
	return wt.ConnectBy(promise, new_conn)
}

func (wt *WSTran) Reconnect(conn fairy.Conn) (fairy.Future, error) {
	if wt.IsStopped() {
		return nil, fmt.Errorf("stopped, cannot reconnect")
	}

	interval := wt.CfgReconnectInterval

	promise := base.NewPromise(conn)

	if interval == 0 {
		wt.ConnectBy(promise, conn)
	} else {
		timer.Start(int64(interval*1000), func(*timer.Timer) {
			wt.ConnectBy(promise, conn)
		})

	}

	return promise, nil
}
