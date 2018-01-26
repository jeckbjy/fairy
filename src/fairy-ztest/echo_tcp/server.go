package echo_tcp

import (
	"fairy"
	"fairy-ztest/msg"
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/tcp"
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
	fmt.Println("Start Server!")

	fairy.RegisterMessage(&msg.EchoMsg{})
	fairy.RegisterHandler(&msg.EchoMsg{}, OnServerEcho)

	transport := tcp.NewTransport()
	transport.AddFilters(
		filter.NewTransportFilter(),
		filter.NewFrameFilter(frame.NewLineFrame()),
		filter.NewPacketFilter(identity.NewStringIdentity(), codec.NewJsonCodec()),
		filter.NewExecutorFilter())

	transport.Listen(":8888", 0)
	transport.Start()

	fairy.WaitExit()
	fmt.Println("Stop Server!")
}
