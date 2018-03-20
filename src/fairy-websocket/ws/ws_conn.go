package ws

import (
	"fairy"
	"fairy/snet"
	"fmt"

	"github.com/gorilla/websocket"
)

func newConn(tran fairy.Transport, side bool, tag interface{}) *snet.StreamConn {
	conn := snet.NewConn(&wsConn{}, tran, side, tag)
	return conn
}

type wsConn struct {
	*websocket.Conn
}

func (wc *wsConn) Open(conn interface{}) {
	wc.Conn = conn.(*websocket.Conn)
}

func (wc *wsConn) Read(reader *fairy.Buffer, cap int) error {
	mtype, data, err := wc.ReadMessage()
	if err != nil {
		return err
	}

	switch mtype {
	case websocket.TextMessage, websocket.BinaryMessage:
		reader.Append(data)
		return nil
	case websocket.CloseMessage:
		return fmt.Errorf("close")
	}

	return nil
}

func (wc *wsConn) Write(buf []byte) error {
	return wc.WriteMessage(websocket.BinaryMessage, buf)
}
