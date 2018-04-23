package tcp

import (
	"net"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/snet"
)

func newConn(tran fairy.Tran, side bool, tag interface{}) *snet.StreamConn {
	conn := snet.NewConn(&zTcpConn{}, tran, side, tag)
	return conn
}

type zTcpConn struct {
	net.Conn
}

func (tc *zTcpConn) Open(conn interface{}) {
	tc.Conn = conn.(net.Conn)
}

func (tc *zTcpConn) Read(reader *fairy.Buffer, cap int) error {
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

func (tc *zTcpConn) Write(buf []byte) error {
	_, err := tc.Conn.Write(buf)
	return err
}
