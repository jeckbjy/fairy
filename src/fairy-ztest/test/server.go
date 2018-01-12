package test

import (
	"fairy"
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/log"
	"fairy/tcp"
	"fairy/timer"
	"fmt"
)

func OnLogin(conn fairy.Connection, packet fairy.Packet) {
	msg := packet.GetMessage().(string)
	fmt.Println(msg)
}

func OnTimeout(t *timer.Timer) {
	log.Debug("OnTimeout!")
}

func StartServer() {
	fmt.Println("Start Server!")

	// 定时器
	// fairy.StartTimer()
	// fairy.StartDelayTimer()
	timer.StartDelay(10, OnTimeout)
	// register
	fairy.RegisterHandler(1, OnLogin)

	transport := tcp.NewTransport()
	transport.AddFilters(
		filter.NewTransportFilter(),
		filter.NewFrameFilter(frame.NewVarintLengthFrame()),
		filter.NewPacketFilter(identity.NewStringIdentity(), codec.NewJsonCodec()),
		filter.NewExecutorFilter())

	transport.Listen(":8888", 0)
	transport.Start(true)

	fmt.Println("Stop Server!")
}
