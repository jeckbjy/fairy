package test

import (
	"fairy"
	"fairy/filter"
	"fairy/tcp"
	"testing"
)

func TestTelnet(t *testing.T) {
	tran := tcp.NewTransport()
	tran.AddFilters(
		filter.NewTransportFilter(),
		filter.NewTelnetFilter())
	tran.Listen(":8080", 0)
	fairy.WaitExit()
}
