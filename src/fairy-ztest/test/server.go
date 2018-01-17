package test

import (
	"fairy"
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/tcp"
	"fmt"
)

type LoginReq struct {
	Account  string
	Password string
}

type LoginRsp struct {
	Error int
}

func OnLogin(conn fairy.Connection, packet fairy.Packet) {
	msg := packet.GetMessage().(string)
	fmt.Println(msg)
}

func OnTimeout(t *fairy.Timer) {
	fairy.Debug("OnTimeout!")
}

func StartServer() {
	fmt.Println("Start Server!")

	// 定时器
	fairy.StartTimer(10, OnTimeout)
	// register
	fairy.RegisterMessage(&LoginReq{})
	fairy.RegisterMessage(&LoginRsp{})
	fairy.RegisterHandler(&LoginReq{}, OnLogin)

	transport := tcp.NewTransport()
	transport.AddFilters(
		filter.NewTransportFilter(),
		filter.NewFrameFilter(frame.NewLineFrame()),
		filter.NewPacketFilter(identity.NewStringIdentity(), codec.NewJsonCodec()),
		filter.NewExecutorFilter())

	transport.Listen(":8888", 0)
	transport.Start(true)

	fmt.Println("Stop Server!")
}
