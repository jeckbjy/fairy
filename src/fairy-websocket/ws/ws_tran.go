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
	new_conn := NewConn(self.owner, true, self.kind)
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

func (wt *WSTran) ConnectBy(promise fairy.Promise, new_conn fairy.Connection) (fairy.Future, error) {
	wt.AddGroup()
	stream_conn := new_conn.(*snet.StreamConn)
	host := stream_conn.GetHost()
	count_max := wt.CfgReconnectCount

	go func() {
		ok := false
		for i := 0; i < count_max && !wt.IsStopped(); i++ {
			// logging ??
			conn, _, err := websocket.DefaultDialer.Dial(host, nil)

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

	new_conn := NewConn(wt, false, kind)
	new_conn.SetHost(host)
	promise := base.NewPromise(new_conn)
	return wt.ConnectBy(promise, new_conn)
}

func (wt *WSTran) Reconnect(conn fairy.Connection) (fairy.Future, error) {
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

// func NewTransport() fairy.Transport {
// 	ws := &WSTransport{}
// 	ws.NewBase()
// 	ws.wg = sync.WaitGroup{}
// 	return ws
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true
// 	},
// }

// type ServeHttpHandler struct {
// 	kind  int
// 	owner *WSTransport
// }

// func (self *ServeHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		return
// 	}

// 	// kind
// 	newConn := NewConnection(self.owner, self.owner.GetFilterChain(), true, self.kind)
// 	newConn.Open(conn)
// }

// type WSTransport struct {
// 	base.Transport
// 	listeners []net.Listener
// 	wg        sync.WaitGroup
// }

// func (self *WSTransport) Listen(host string, kind int) error {
// 	listener, err := net.Listen("tcp", host)
// 	if err != nil {
// 		return err
// 	}

// 	self.wg.Add(1)
// 	self.listeners = append(self.listeners, listener)
// 	go func() {
// 		defer self.wg.Done()
// 		for {
// 			svr := http.Server{Handler: &ServeHttpHandler{kind: kind, owner: self}}
// 			err := svr.Serve(listener)
// 			if err != nil {
// 				break
// 			}
// 		}
// 	}()

// 	return nil
// }

// func (self *WSTransport) Connect(host string, kind int) (fairy.ConnectFuture, error) {
// 	// convert url:localhost:8888->ws://localhost:8888
// 	pos := strings.Index(host, "//")
// 	if pos == -1 {
// 		host = "ws://" + host
// 	}

// 	// u, err := url.Parse(host)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// if (u.Scheme != "ws") && (u.Scheme != "wss") {
// 	// 	return nil, fmt.Errorf("malformed ws or wss URL")
// 	// }

// 	newConn := NewConnection(self, self.GetFilterChain(), false, kind)
// 	newConn.Host = host
// 	future := base.NewConnectFuture(newConn)
// 	self.ConnectBy(future, newConn, host)
// 	return future, nil
// }

// func (self *WSTransport) ConnectBy(future fairy.ConnectFuture, newConn *WSConnection, host string) {
// 	self.wg.Add(1)
// 	go func() {
// 		// wait for close
// 		defer self.wg.Done()

// 		conn, _, err := websocket.DefaultDialer.Dial(host, nil)
// 		if future == nil || future.Result() != fairy.FUTURE_RESULT_TIMEOUT {
// 			result := 0
// 			if err == nil {
// 				newConn.Open(conn)
// 				result = fairy.FUTURE_RESULT_SUCCEED
// 			} else {
// 				fairy.Debug("connect fail!:%+v", err)
// 				newConn.HandleError(newConn, fairy.ErrConnectFail)
// 				result = fairy.FUTURE_RESULT_FAIL
// 			}

// 			if future != nil {
// 				future.Done(result)
// 			}
// 		} else if err != nil {
// 			conn.Close()
// 		}
// 	}()
// }

// func (t *WSTransport) TryReconnect(conn *WSConnection) bool {
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

// func (t *WSTransport) Stop() {
// 	// close all listener
// 	for _, listener := range t.listeners {
// 		listener.Close()
// 	}

// 	t.listeners = nil

// 	// stop reconnect
// 	t.CfgReconnectInterval = -1
// }

// func (t *WSTransport) Wait() {
// 	t.wg.Wait()
// }

// func (t *WSTransport) OnExit() {
// 	t.Stop()
// }
