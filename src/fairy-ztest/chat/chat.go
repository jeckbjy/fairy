package chat

import (
	"fairy"
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/log"
	"fairy/tcp"
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
		filter.NewLog(),
		filter.NewFrame(frame.NewLine()),
		filter.NewPacket(identity.NewString(), codec.NewJson()),
		filter.NewExecutor())

	// step4: listen or connect
	tran.Listen(":8080", 0)
	// step5: wait finish
	fairy.WaitExit()
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

	// step3: create transport and add filters
	tran := tcp.NewTransport()
	tran.AddFilters(
		filter.NewLog(),
		filter.NewFrame(frame.NewLine()),
		filter.NewPacket(identity.NewString(), codec.NewJson()),
		filter.NewExecutor())

	tran.AddFilters(filter.NewConnect(func(conn fairy.Conn) {
		// send msg to server
		req := &ChatMsg{}
		req.Content = "hello word!"
		conn.Send(req)
	}))

	// step4: listen or connect
	tran.Connect("localhost:8080", 0)
	// step5: wait finish
	fairy.WaitExit()
}
