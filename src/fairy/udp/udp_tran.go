package udp

import (
	"fairy"
	"fairy/base"
	"fmt"
	"net"
)

type UdpTran struct {
	base.Tran
}

func (ut *UdpTran) Listen(host string, kind int) error {
	udpaddr, err := net.ResolveUDPAddr("udp", host)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		return err
	}

	fmt.Printf("%+v", conn)

	return nil
}

func (ut *UdpTran) Connect(host string, kind int) (fairy.Future, error) {
	return nil, nil
}

func (ut *UdpTran) Reconnect(conn fairy.Conn) (fairy.Future, error) {
	return nil, fmt.Errorf("udp transport cannot reconnect")
}
