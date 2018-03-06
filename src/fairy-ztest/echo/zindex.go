package echo

import (
	"flag"
)

func Test() {
	// example: ./test -s server -n tcp -m json
	pside := flag.String("s", "client", "test mode:server or client")
	pnetmode := flag.String("n", "tcp", "network mode:tcp,ws,kcp")
	pmsgmode := flag.String("m", "json", "proto mode:json,pb")
	flag.Parse()

	if *pside == "server" {
		StartServer(*pnetmode, *pmsgmode)
	} else {
		StartClient(*pnetmode, *pmsgmode)
	}
}
