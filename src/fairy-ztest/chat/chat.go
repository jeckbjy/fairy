package chat

import (
	"fairy"
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/log"
	"fairy/tcp"
	"fairy/timer"
)

type ChatMsg struct {
	Content string
}

func StartServer() {
	log.Debug("start server")
	// step1: register message
	fairy.RegisterMessage(&ChatMsg{})

	// step2: register handler
	fairy.RegisterHandler(&ChatMsg{}, func(conn fairy.Conn, pkt fairy.Packet) {
		req := pkt.GetMessage().(*ChatMsg)
		log.Debug("client msg:%+v", req.Content)

		rsp := &ChatMsg{}
		rsp.Content = "welcome boy!"
		conn.Send(rsp)
	})

	// step3: create transport and add filters
	tran := tcp.NewTransport()
	tran.AddFilters(
		filter.NewLoggingEx(filter.LoggingFilterAll),
		filter.NewFrame(frame.NewLine()),
		filter.NewPacket(identity.NewString(), codec.NewJson()),
		filter.NewExecutor())

	// step4: listen or connect
	tran.Listen(":8080", 0)
}

func StartClient() {
	log.Debug("start client")
	// step1: register message
	fairy.RegisterMessage(&ChatMsg{})

	// step2: register handler
	fairy.RegisterHandler(&ChatMsg{}, func(conn fairy.Conn, pkt fairy.Packet) {
		req := pkt.GetMessage().(*ChatMsg)
		log.Debug("server msg:%+v", req.Content)
	})

	var gConn fairy.Conn
	// step3: create transport and add filters
	tran := tcp.NewTransport()
	tran.AddFilters(
		filter.NewLoggingEx(filter.LoggingFilterAll),
		filter.NewFrame(frame.NewLine()),
		filter.NewPacket(identity.NewString(), codec.NewJson()),
		filter.NewExecutor())

	tran.AddFilters(filter.NewConnect(func(conn fairy.Conn) {
		// send msg to server
		req := &ChatMsg{}
		req.Content = "hello word!"
		conn.Send(req)
		gConn = conn
	}))

	// add timer for send message
	timer.Start(1000, func(t *timer.Timer) {
		log.Debug("Ontimeout")
		req := &ChatMsg{}
		req.Content = "hello word!"
		gConn.Send(req)
		t.Restart()
	})

	// step4: listen or connect
	tran.Connect("localhost:8080", 0)
}
