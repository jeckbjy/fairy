package echo

import (
	"fairy"
	"fairy-ztest/echo/json"
	"fairy-ztest/echo/pb"
	"fairy/filter"
	"fairy/log"
	"fairy/timer"
	"fairy/util"
	"fmt"
)

var gClient fairy.Conn

func SendEchoToServer() {
	log.Debug("send msg to server!")
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

func OnTimeout(t *timer.Timer) {
	log.Debug("timeout")
	if gClient == nil {
		return
	}

	SendEchoToServer()
	t.Restart()
}

func OnConnected(conn fairy.Conn) {
	log.Debug("OnConnected")
	gClient = conn
	SendEchoToServer()
}

func OnClientEcho(conn fairy.Conn, packet fairy.Packet) {
	rsp := packet.GetMessage()
	log.Debug("Recv server echo: %+v", rsp)
}

func StartClient(net_mode string, msg_mode string) {
	fmt.Printf("start client:net_mode=%v, msg_mode=%v\n", net_mode, msg_mode)

	SetMsgMode(msg_mode)
	RegisterMsg(msg_mode, OnClientEcho)

	tran := NewTransport(net_mode, msg_mode)
	tran.AddFilters(filter.NewConnect(OnConnected))

	tran.Connect("localhost:8888", 0)
	tran.Start()

	timer.Start(500, OnTimeout)
	fairy.WaitExit()
	fmt.Printf("stop client!\n")
}
