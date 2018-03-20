package tcp

import (
	"fairy"
	"fairy/snet"
	"net"
)

func newConn(tran fairy.Transport, side bool, tag interface{}) *snet.StreamConn {
	conn := snet.NewConn(&TcpConn{}, tran, side, tag)
	return conn
}

type TcpConn struct {
	net.Conn
}

func (tc *TcpConn) Open(conn interface{}) {
	tc.Conn = conn.(net.Conn)
}

func (tc *TcpConn) Read(reader *fairy.Buffer, cap int) error {
	data := reader.GetSpace()
	hasSpace := true
	if data == nil {
		data = make([]byte, cap)
		hasSpace = false
	}

	n, err := tc.Conn.Read(data)
	if err != nil {
		return err
	}

	if hasSpace {
		reader.ExtendSpace(n)
	} else {
		reader.Append(data[:n])
	}

	return nil
}

func (tc *TcpConn) Write(buf []byte) error {
	_, err := tc.Conn.Write(buf)
	return err
}
