package ws

import (
	"fairy"
	"fairy/base"
	"net"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

func NewTransport() fairy.Transport {
	ws := &WSTransport{}
	ws.NewBase()
	ws.stopped = make(chan bool)
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
	stopped chan bool
	wg      sync.WaitGroup
}

func (self *WSTransport) Listen(host string, kind int) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
		return
	}

	self.wg.Add(1)
	go func() {
		defer self.wg.Done()

		for {
			select {
			case <-self.stopped:
				break
			default:
				svr := http.Server{Handler: &ServeHttpHandler{kind: kind, owner: self}}
				err := svr.Serve(listener)
				if err != nil {
					fairy.Error("%+v", err)
				}
			}
		}
	}()
}

func (self *WSTransport) Connect(host string, kind int) fairy.ConnectFuture {
	// url check
	u, err := url.Parse(host)
	if err != nil {
		return nil
	}

	if (u.Scheme != "ws") && (u.Scheme != "wss") {
		return nil
	}

	newConn := NewConnection(self, self.GetFilterChain(), false, kind)
	future := base.NewConnectFuture(newConn)
	self.ConnectBy(future, newConn, host)
	return future
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

func (self *WSTransport) Start() {
	self.wg.Add(1)
}

func (self *WSTransport) Stop() {
	close(self.stopped)
	self.wg.Done()
	self.wg.Wait()
}

func (self *WSTransport) Wait() {
	self.wg.Wait()
}

func (self *WSTransport) OnExit() {
	self.Stop()
}

func (t *WSTransport) Reconnect(conn *WSConnection) {
	// 断线重连
	if t.CfgReconnectInterval == 0 {
		t.ConnectBy(nil, conn, conn.Host)
	} else {
		fairy.StartTimer(int64(t.CfgReconnectInterval*1000), func(*fairy.Timer) {
			t.ConnectBy(nil, conn, conn.Host)
		})
	}
}

func (t *WSTransport) HandleConnClose(conn *WSConnection) {
	if conn.IsClientSide() && t.IsNeedReconnect() {
		t.Reconnect(conn)
	} else {
		// reomve conn??
	}
}
