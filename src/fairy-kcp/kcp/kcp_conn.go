package kcp

import (
	"fairy"
	"fairy/snet"
	"net"
)

func NewConn(tran fairy.Transport, side bool, kind int) *snet.StreamConn {
	conn := snet.NewConn(&KcpConn{}, tran, side, kind)
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
