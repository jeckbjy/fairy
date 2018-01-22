package websocket

import (
	"fairy"
	"fairy/base"
	"net"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

func NewTransport() fairy.Transport {
	ws := &WSTransport{}
	ws.BaseTransport.New()
	ws.stopFlag = make(chan bool)
	ws.waitGroup = sync.WaitGroup{}
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
	stopFlag  chan bool
	waitGroup sync.WaitGroup
}

func (self *WSTransport) Listen(host string, kind int) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
		return
	}

	self.waitGroup.Add(1)
	go func() {
		defer self.waitGroup.Done()

		for {
			select {
			case <-self.stopFlag:
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
	newConn := NewConnection(self, self.GetFilterChain(), false, kind)
	future := base.NewConnectFuture(newConn)
	self.ConnectBy(future, newConn, host)
	return future
}

func (self *WSTransport) ConnectBy(future fairy.ConnectFuture, newConn *WSConnection, host string) {
	self.waitGroup.Add(1)
	go func() {
		// wait for close
		defer self.waitGroup.Done()

		conn, _, err := websocket.DefaultDialer.Dial(host, nil)
		if future == nil || future.Result() != fairy.FUTURE_RESULT_TIMEOUT {
			result := 0
			if err != nil {
				newConn.Open(conn)
				newConn.HandleOpen(newConn)
				result = fairy.FUTURE_RESULT_SUCCEED
			} else {
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
	self.waitGroup.Add(1)
}

func (self *WSTransport) Stop() {
	close(self.stopFlag)
	self.waitGroup.Done()
	self.waitGroup.Wait()
}

func (self *WSTransport) Wait() {
	self.waitGroup.Wait()
}

func (self *WSTransport) OnExit() {
	self.Stop()
}
