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

func OnServerEcho(conn fairy.Connection, packet fairy.Packet) {
	req := packet.GetMessage().(*msg.EchoMsg)
	fairy.Debug(" OnServerEcho: %+v", req)
	rsp := &msg.EchoMsg{}
	rsp.Info = "server rsp!"
	rsp.Timestamp = util.Now()
	conn.Send(rsp)
}

func StartServer() {
	fmt.Println("Start WebSocket Server!")

	fairy.RegisterMessage(&msg.EchoMsg{})
	fairy.RegisterHandler(&msg.EchoMsg{}, OnServerEcho)

	transport := ws.NewTransport()
	transport.AddFilters(
		filter.NewTransportFilter(),
		filter.NewFrameFilter(frame.NewLineFrame()),
		filter.NewPacketFilter(identity.NewStringIdentity(), codec.NewJsonCodec()),
		filter.NewExecutorFilter())

	transport.Listen(":8080", 0)
	transport.Start()

	fairy.WaitExit()
	fmt.Println("Stop Server!")
}
