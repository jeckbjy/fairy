package chat

import (
	"fairy"
	"flag"
)

func Test() {
	// example: ./test -s server -n tcp -m json
	pside := flag.String("s", "server", "test mode:server or client")
	flag.Parse()

	if *pside == "server" {
		StartServer()
	} else {
		StartClient()
	}

	fairy.WaitExit()
}
