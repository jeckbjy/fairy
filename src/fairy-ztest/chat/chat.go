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
	"fairy/util"
)

type ChatMsg struct {
	Content   string
	Timestamp int64
}

func StartServer() {
	log.Debug("start server")
	// step1: register message
	fairy.RegisterMessage(&ChatMsg{})

	// step2: register handler
	fairy.RegisterHandler(&ChatMsg{}, func(conn fairy.Conn, pkt fairy.Packet) {
		req := pkt.GetMessage().(*ChatMsg)
		log.Debug("client msg:%+v", req)

		rsp := &ChatMsg{}
		rsp.Content = "welcome boy!"
		rsp.Timestamp = util.Now()
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
		log.Debug("server msg:%+v", req)
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
		req.Timestamp = util.Now()
		gConn.Send(req)
		if gConn.IsActive() {
			t.Restart()
		}
	})

	// step4: listen or connect
	tran.Connect("localhost:8080", 0)
}
