package echo_ws

import (
	"fairy"
	"fairy-websocket/ws"
	"fairy-ztest/msg"
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/util"
	"fmt"
)

var gClient fairy.Connection

func OnTimeout(timer *fairy.Timer) {
	if gClient == nil {
		return
	}

	req := &msg.EchoMsg{}
	req.Info = "Client Echo!"
	req.Timestamp = util.Now()
	gClient.Send(req)
}

func OnConnected(conn fairy.Connection) {
	fairy.Debug("OnConnected")
	gClient = conn

	req := &msg.EchoMsg{}
	req.Info = "Client Echo!"
	req.Timestamp = util.Now()
	gClient.Send(req)
}

func OnClientEcho(conn fairy.Connection, packet fairy.Packet) {
	rsp := packet.GetMessage().(*msg.EchoMsg)
	fairy.Debug("server echo:%+v", rsp)

	req := &msg.EchoMsg{}
	req.Info = "Client Echo!"
	req.Timestamp = util.Now()
	gClient.Send(req)
}

func StartClient() {
	fmt.Println("Start WebSocket Client!")

	fairy.RegisterMessage(&msg.EchoMsg{})
	fairy.RegisterHandler(&msg.EchoMsg{}, OnClientEcho)

	// fairy.StartTimer(10000, OnTimeout)

	transport := ws.NewTransport()
	transport.AddFilters(
		filter.NewTransportFilter(),
		filter.NewFrameFilter(frame.NewLineFrame()),
		filter.NewPacketFilter(identity.NewStringIdentity(), codec.NewJsonCodec()),
		filter.NewExecutorFilter(),
		filter.NewConnectFilter(OnConnected))

	transport.Connect("ws://localhost:8080", 0)
	transport.Start()
	fairy.WaitExit()
	fmt.Println("Stop Client!")
}
