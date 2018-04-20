package soa

import (
	"github.com/jeckbjy/fairy"
)

func register(msg interface{}, cb fairy.HandlerCB) {
	fairy.RegisterMessage(msg)
	fairy.RegisterHandler(msg, cb)
}
