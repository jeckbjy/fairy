package test

import (
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
		filter.NewFrameFilter(frame.NewVarintLengthFrame()),
		filter.NewPacketFilter(identity.NewStringIdentity(), codec.NewJsonCodec()))

	transport.Connect("127.0.0.1:8888", 0)
	transport.Start(true)
	fmt.Println("Stop Client!")
}
