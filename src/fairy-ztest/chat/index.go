package chat

import "flag"

func Test() {
	// example: ./test -s server -n tcp -m json
	pside := flag.String("s", "client", "test mode:server or client")
	flag.Parse()

	if *pside == "server" {
		StartServer()
	} else {
		StartClient()
	}
}
