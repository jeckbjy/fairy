package tcp

import (
	"fairy"
	"fairy/snet"
	"net"
)

func newConn(tran fairy.Transport, side bool, kind int) *snet.StreamConn {
	conn := snet.NewConn(&TcpConn{}, tran, side, kind)
	return conn
}

type TcpConn struct {
	net.Conn
}

func (tc *TcpConn) Open(conn interface{}) {
	tc.Conn = conn.(net.Conn)
}

func (tc *TcpConn) Read(cap int) ([]byte, error) {
	data := make([]byte, cap)
	n, err := tc.Conn.Read(data)
	if err != nil {
		return nil, err
	}

	return data[:n], nil
}

func (tc *TcpConn) Write(buf []byte) error {
	_, err := tc.Conn.Write(buf)
	return err
}
