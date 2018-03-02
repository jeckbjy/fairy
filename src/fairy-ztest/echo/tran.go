package echo

import (
	"fairy"
	"fairy-kcp/kcp"
	"fairy-protobuf/pbcodec"
	"fairy-websocket/ws"
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/tcp"
)

func NewTransport(net_mode string, msg_mode string) fairy.Transport {
	var tran fairy.Transport
	switch net_mode {
	case "ws":
		tran = ws.NewTransport()
	case "kcp":
		tran = kcp.NewTransport()
	default:
		// tcp
		tran = tcp.NewTransport()
	}

	var zframe fairy.Frame
	var zidentity fairy.Identity
	var zcodec fairy.Codec

	switch msg_mode {
	case "pb":
		zframe = frame.NewVarintLengthFrame()
		zidentity = identity.NewIntegerIdentity()
		zcodec = pbcodec.New()
	default:
		// json
		zframe = frame.NewLineFrame()
		zidentity = identity.NewStringIdentity()
		zcodec = codec.NewJson()
	}

	tran.AddFilters(
		filter.NewLogFilter(),
		filter.NewFrameFilter(zframe),
		filter.NewPacketFilter(zidentity, zcodec),
		filter.NewExecutorFilter())

	return tran
}
