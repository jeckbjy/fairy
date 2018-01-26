package main

import (
	"fairy-ztest/echo_tcp"
	"fairy-ztest/echo_ws"
	"flag"
)

func main() {
	mode_ptr := flag.String("m", "client-ws", "server mode")
	flag.Parse()
	mode := *mode_ptr

	switch mode {
	case "server":
		echo_tcp.StartServer()
	case "client":
		echo_tcp.StartClient()
	case "server-ws":
		echo_ws.StartServer()
	case "client-ws":
		echo_ws.StartClient()
	}
}
