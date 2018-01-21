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

func StartClient() {
	fmt.Println("Start Client!")

	transport := tcp.NewTransport()
	transport.AddFilters(
		filter.NewTransportFilter(),
		filter.NewFrameFilter(frame.NewLineFrame()),
		filter.NewPacketFilter(identity.NewStringIdentity(), codec.NewJsonCodec()),
		filter.NewExecutorFilter())

	transport.Connect("127.0.0.1:8888", 0)
	transport.Start()
	fairy.WaitExit()
	fmt.Println("Stop Client!")
}
