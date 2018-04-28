package ws

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
	"github.com/jeckbjy/fairy/snet"
	"github.com/jeckbjy/fairy/timer"
)

// NewTran create websocket transport
func NewTran() fairy.Tran {
	tran := &wsTran{}
	tran.Create()
	return tran
}

type wsServeHTTPHandler struct {
	tag   interface{}
	owner *wsTran
}

func (self *wsServeHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := self.owner.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	new_conn := newConn(self.owner, true, self.tag)
	new_conn.Open(conn)
}

type wsTran struct {
	snet.StreamTran
	websocket.Upgrader
}

func (wt *wsTran) Create() {
	wt.StreamTran.Create()
	wt.Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func (wt *wsTran) Listen(host string, tag interface{}) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	wt.AddListener(listener)

	wt.AddGroup()
	go func() {
		for {
			svr := http.Server{Handler: &wsServeHTTPHandler{tag: tag, owner: wt}}
			err := svr.Serve(listener)
			if err != nil {
				break
			}
		}

		wt.Done()
	}()

	return nil
}

func (wt *wsTran) ConnectBy(promise fairy.Promise, new_conn fairy.Conn) (fairy.Future, error) {
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

func (wt *wsTran) Connect(host string, tag interface{}) (fairy.Future, error) {
	// convert url:localhost:8888->ws://localhost:8888
	pos := strings.Index(host, "//")
	if pos == -1 {
		host = "ws://" + host
	}

	new_conn := newConn(wt, false, tag)
	new_conn.SetHost(host)
	promise := base.NewPromise(new_conn)
	return wt.ConnectBy(promise, new_conn)
}

func (wt *wsTran) Reconnect(conn fairy.Conn) (fairy.Future, error) {
	if wt.IsStopped() {
		return nil, fmt.Errorf("stopped, cannot reconnect")
	}

	interval := wt.CfgReconnectInterval

	promise := base.NewPromise(conn)

	if interval == 0 {
		wt.ConnectBy(promise, conn)
	} else {
		timer.Start(timer.ModeDelay, int64(interval*1000), func() {
			wt.ConnectBy(promise, conn)
		})

	}

	return promise, nil
}
