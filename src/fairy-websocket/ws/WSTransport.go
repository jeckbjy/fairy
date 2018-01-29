package ws

import (
	"fairy"
	"fairy/base"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

func NewTransport() fairy.Transport {
	ws := &WSTransport{}
	ws.NewBase()
	ws.wg = sync.WaitGroup{}
	return ws
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ServeHttpHandler struct {
	kind  int
	owner *WSTransport
}

func (self *ServeHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// kind
	newConn := NewConnection(self.owner, self.owner.GetFilterChain(), true, self.kind)
	newConn.Open(conn)
	newConn.HandleOpen(newConn)
}

type WSTransport struct {
	base.BaseTransport
	wg        sync.WaitGroup
	listeners []net.Listener
}

func (self *WSTransport) Listen(host string, kind int) error {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	self.wg.Add(1)
	self.listeners = append(self.listeners, listener)
	go func() {
		defer self.wg.Done()
		for {
			svr := http.Server{Handler: &ServeHttpHandler{kind: kind, owner: self}}
			err := svr.Serve(listener)
			if err != nil {
				break
			}
		}
	}()

	return nil
}

func (self *WSTransport) Connect(host string, kind int) (fairy.ConnectFuture, error) {
	// url check
	u, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	if (u.Scheme != "ws") && (u.Scheme != "wss") {
		return nil, fmt.Errorf("malformed ws or wss URL")
	}

	newConn := NewConnection(self, self.GetFilterChain(), false, kind)
	newConn.Host = host
	future := base.NewConnectFuture(newConn)
	self.ConnectBy(future, newConn, host)
	return future, nil
}

func (self *WSTransport) ConnectBy(future fairy.ConnectFuture, newConn *WSConnection, host string) {
	self.wg.Add(1)
	go func() {
		// wait for close
		defer self.wg.Done()

		conn, _, err := websocket.DefaultDialer.Dial(host, nil)
		if future == nil || future.Result() != fairy.FUTURE_RESULT_TIMEOUT {
			result := 0
			if err == nil {
				newConn.Open(conn)
				newConn.HandleOpen(newConn)
				result = fairy.FUTURE_RESULT_SUCCEED
			} else {
				fairy.Debug("connect fail!:%+v", err)
				newConn.HandleError(newConn, fairy.ErrConnectFail)
				result = fairy.FUTURE_RESULT_FAIL
			}

			if future != nil {
				future.Done(result)
			}
		} else if err != nil {
			conn.Close()
		}
	}()
}

func (t *WSTransport) TryReconnect(conn *WSConnection) bool {
	if conn.IsClientSide() && t.IsNeedReconnect() {
		// 断线重连
		if t.CfgReconnectInterval == 0 {
			t.ConnectBy(nil, conn, conn.Host)
		} else {
			fairy.StartTimer(int64(t.CfgReconnectInterval*1000), func(*fairy.Timer) {
				t.ConnectBy(nil, conn, conn.Host)
			})
		}

		return true
	}

	return false
}

func (t *WSTransport) Stop() {
	// close all listener
	for _, listener := range t.listeners {
		listener.Close()
	}

	t.listeners = nil

	// stop reconnect
	t.CfgReconnectInterval = -1
}

func (t *WSTransport) Wait() {
	t.wg.Wait()
}

func (t *WSTransport) OnExit() {
	t.Stop()
}
