package kcp

import (
	"net"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/snet"
)

func NewConn(tran fairy.Tran, side bool, tag interface{}) *snet.StreamConn {
	conn := snet.NewConn(&KcpConn{}, tran, side, tag)
	return conn
}

type KcpConn struct {
	net.Conn
}

func (kc *KcpConn) Open(conn interface{}) {
	kc.Conn = conn.(net.Conn)
}

func (kc *KcpConn) Read(reader *fairy.Buffer, cap int) error {
	data := reader.GetSpace()
	hasSpace := true
	if data == nil {
		data = make([]byte, cap)
		hasSpace = false
	}

	n, err := kc.Conn.Read(data)
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

func (kc *KcpConn) Write(buf []byte) error {
	_, err := kc.Conn.Write(buf)
	return err
}
