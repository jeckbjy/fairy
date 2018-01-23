package echo

import (
	"fairy"
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/tcp"
	"fairy/util"
	"fmt"
)

var gClient fairy.Connection

func OnTimeout(timer *fairy.Timer) {
	if gClient == nil {
		return
	}

	req := &EchoMsg{}
	req.Info = "Client Echo!"
	req.Timestamp = util.Now()
	gClient.Send(req)
}

func OnConnected(conn fairy.Connection) {
	fairy.Debug("OnConnected")
	gClient = conn

	req := &EchoMsg{}
	req.Info = "Client Echo!"
	req.Timestamp = util.Now()
	gClient.Send(req)
}

func OnClientEcho(conn fairy.Connection, packet fairy.Packet) {
	rsp := packet.GetMessage().(*EchoMsg)
	fairy.Debug("%+v", rsp)
}

func StartClient() {
	fmt.Println("Start Client!")

	fairy.RegisterMessage(&EchoMsg{})
	fairy.RegisterHandler(&EchoMsg{}, OnClientEcho)

	// fairy.StartTimer(10000, OnTimeout)

	transport := tcp.NewTransport()
	transport.AddFilters(
		filter.NewTransportFilter(),
		filter.NewFrameFilter(frame.NewLineFrame()),
		filter.NewPacketFilter(identity.NewStringIdentity(), codec.NewJsonCodec()),
		filter.NewExecutorFilter(),
		filter.NewConnectFilter(OnConnected))

	transport.Connect("127.0.0.1:8888", 0)
	transport.Start()
	fairy.WaitExit()
	fmt.Println("Stop Client!")
}
