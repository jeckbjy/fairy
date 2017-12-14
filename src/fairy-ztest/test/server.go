package test

import (
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/tcp"
	"fmt"
)

func StartServer() {
	fmt.Println("Start Server!")

	transport := tcp.NewTransport()
	transport.AddFilters(
		filter.NewTransportFilter(),
		filter.NewFrameFilter(frame.NewVarintLengthFrame()),
		filter.NewPacketFilter(identity.NewDefaultStringIdentity(), codec.NewJsonCodec()))

	transport.Listen(":8888", 0)
	transport.Start(true)

	fmt.Println("Stop Server!")
}
