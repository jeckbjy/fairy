package echo

import (
	"fairy"
	"fairy/filter"
	"fairy/util"
	"fmt"
)

var gClient fairy.Connection

func SendEchoToServer() {
	req := &EchoMsg{}
	req.Info = "Client Echo!"
	req.Timestamp = util.Now()
	gClient.Send(req)
}

func OnTimeout(timer *fairy.Timer) {
	if gClient == nil {
		return
	}

	SendEchoToServer()
	timer.Restart(1000)
}

func OnConnected(conn fairy.Connection) {
	fairy.Debug("OnConnected")
	gClient = conn
	SendEchoToServer()
}

func OnClientEcho(conn fairy.Connection, packet fairy.Packet) {
	rsp := packet.GetMessage().(*EchoMsg)
	fairy.Debug("Recv server echo: %+v", rsp)

	// req := &msg.EchoMsg{}
	// req.Info = "Client Echo!"
	// req.Timestamp = util.Now()
	// gClient.Send(req)
}

func StartClient(net_mode string, msg_mode string) {
	fmt.Printf("start client:net_mode=%v, msg_mode=%v\n", net_mode, msg_mode)

	switch msg_mode {
	case "pb":
	default:
		// json
		fairy.RegisterMessage(&EchoMsg{})
		fairy.RegisterHandler(&EchoMsg{}, OnClientEcho)
	}

	var host string
	switch net_mode {
	case "ws":
		host = "ws://localhost:8888"
	default:
		host = "localhost:8888"
	}

	fairy.StartTimer(10000, OnTimeout)

	tran := NewTransport(net_mode, msg_mode)
	tran.AddFilters(filter.NewConnectFilter(OnConnected))

	tran.Connect(host, 0)
	tran.Start()

	fairy.WaitExit()
	fmt.Printf("stop client!\n")
}
