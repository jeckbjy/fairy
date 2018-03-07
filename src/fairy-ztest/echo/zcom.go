package echo

import (
	"fairy"
	"fairy-kcp/kcp"
	"fairy-protobuf/pbcodec"
	"fairy-websocket/ws"
	"fairy-ztest/echo/json"
	"fairy-ztest/echo/pb"
	"fairy/codec"
	"fairy/filter"
	"fairy/frame"
	"fairy/identity"
	"fairy/tcp"
)

var gMsgMode string

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
		zframe = frame.NewVarintLength()
		zidentity = identity.NewInteger()
		zcodec = pbcodec.New()
	default:
		// json
		zframe = frame.NewLine()
		zidentity = identity.NewString()
		zcodec = codec.NewJson()
	}

	tran.AddFilters(
		filter.NewLogging(),
		filter.NewFrame(zframe),
		filter.NewPacket(zidentity, zcodec),
		filter.NewExecutor())

	return tran
}

func RegisterMsg(msg_mode string, cb fairy.HandlerCB) {
	switch msg_mode {
	case "pb":
		// protobuf
		Register(cb, &pb.EchoMsg{}, 1)
	default:
		// json
		Register(cb, &json.EchoMsg{}, 0)
	}
}

func Register(cb fairy.HandlerCB, msg interface{}, id int) {
	if id == 0 {
		fairy.RegisterMessage(msg)
		fairy.RegisterHandler(msg, cb)
	} else {
		fairy.RegisterMessage(msg, id)
		fairy.RegisterHandler(id, cb)
	}

}

func SetMsgMode(mode string) {
	gMsgMode = mode
}

func IsJsonMode() bool {
	return gMsgMode == "json"
}

func IsProtobufMode() bool {
	return gMsgMode == "pb"
}
