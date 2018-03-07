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

func (kc *KcpConn) Read(cap int) ([]byte, error) {
	data := make([]byte, cap)
	n, err := kc.Conn.Read(data)
	if err != nil {
		return nil, err
	}

	return data[:n], nil
}

func (kc *KcpConn) Write(buf []byte) error {
	_, err := kc.Conn.Write(buf)
	return err
}
