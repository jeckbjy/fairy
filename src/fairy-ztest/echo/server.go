package echo

import (
	"fairy"
	"fairy-ztest/echo/json"
	"fairy-ztest/echo/pb"
	"fairy/log"
	"fairy/util"
	"fmt"
)

func OnServerEcho(conn fairy.Conn, packet fairy.Packet) {
	if IsJsonMode() {
		req := packet.GetMessage().(*json.EchoMsg)
		log.Debug("Recv client echo: %+v", req)
		rsp := &json.EchoMsg{}
		rsp.Info = "server rsp!"
		rsp.Timestamp = util.Now()
		conn.Send(rsp)
	} else if IsProtobufMode() {
		req := packet.GetMessage().(*pb.EchoMsg)
		log.Debug("Recv client echo: %+v", req)
		rsp := &pb.EchoMsg{}
		rsp.Info = "server rsp!"
		rsp.Timestamp = util.Now()
		conn.Send(rsp)
	}
}

func StartServer(net_mode string, msg_mode string) {
	fmt.Printf("start server:net_mode=%v, msg_mode=%v\n", net_mode, msg_mode)

	SetMsgMode(msg_mode)
	RegisterMsg(msg_mode, OnServerEcho)
	tran := NewTransport(net_mode, msg_mode)

	tran.Listen(":8888", 0)
	tran.Start()

	fairy.WaitExit()
	fmt.Sprintf("stop server!\n")
}
