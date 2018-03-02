package echo

import (
	"fairy"
	"fairy/util"
	"fmt"
)

func OnServerEcho(conn fairy.Connection, packet fairy.Packet) {
	req := packet.GetMessage().(*EchoMsg)
	fairy.Debug("Recv client echo: %+v", req)
	rsp := &EchoMsg{}
	rsp.Info = "server rsp!"
	rsp.Timestamp = util.Now()
	conn.Send(rsp)
}

func StartServer(net_mode string, msg_mode string) {
	fmt.Printf("start server:net_mode=%v, msg_mode=%v\n", net_mode, msg_mode)

	switch msg_mode {
	case "pb":
	default:
		// json
		fairy.RegisterMessage(&EchoMsg{})
		fairy.RegisterHandler(&EchoMsg{}, OnServerEcho)
	}

	tran := NewTransport(net_mode, msg_mode)

	tran.Listen(":8888", 0)
	tran.Start()

	fairy.WaitExit()
	fmt.Sprintf("stop server!\n")
}
