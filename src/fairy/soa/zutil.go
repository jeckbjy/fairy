package soa

import (
	"fairy"
)

func register(msg interface{}, cb fairy.HandlerCB) {
	fairy.RegisterMessage(msg)
	fairy.RegisterHandler(msg, cb)
}
