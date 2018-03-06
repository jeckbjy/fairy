package echo

import (
	"fairy"
	"fairy-ztest/echo/json"
	"fairy-ztest/echo/pb"
	"fairy/filter"
	"fairy/log"
	"fairy/util"
	"fmt"
)

var gClient fairy.Connection

func SendEchoToServer() {
	if IsJsonMode() {
		req := &json.EchoMsg{}
		req.Info = "Client json.Echo!"
		req.Timestamp = util.Now()
		gClient.Send(req)
	} else {
		req := &pb.EchoMsg{}
		req.Info = "Client pb.Echo!"
		req.Timestamp = util.Now()
		gClient.Send(req)
	}
}

func OnTimeout(timer *fairy.Timer) {
	if gClient == nil {
		return
	}

	SendEchoToServer()
	timer.Restart(1000)
}

func OnConnected(conn fairy.Connection) {
	log.Debug("OnConnected")
	gClient = conn
	SendEchoToServer()
}

func OnClientEcho(conn fairy.Connection, packet fairy.Packet) {
	rsp := packet.GetMessage()
	log.Debug("Recv server echo: %+v", rsp)

	// req := &msg.EchoMsg{}
	// req.Info = "Client Echo!"
	// req.Timestamp = util.Now()
	// gClient.Send(req)
}

func StartClient(net_mode string, msg_mode string) {
	fmt.Printf("start client:net_mode=%v, msg_mode=%v\n", net_mode, msg_mode)

	SetMsgMode(msg_mode)
	RegisterMsg(msg_mode, OnClientEcho)

	tran := NewTransport(net_mode, msg_mode)
	tran.AddFilters(filter.NewConnectFilter(OnConnected))

	tran.Connect("localhost:8888", 0)
	tran.Start()

	fairy.StartTimer(10000, OnTimeout)
	fairy.WaitExit()
	fmt.Printf("stop client!\n")
}
