package main

import (
	"bufio"
	"fairy-ztest/echo_tcp"
	"fairy-ztest/echo_ws"
	"flag"
	"fmt"
	"net"
	"strings"
)

func TestFairy() {
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

func TestTelnet() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen err")
	}
	fmt.Println("listen post 8888!")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept err!")
			break
		}
		go func() {
			for {
				conn.Write([]byte("fairy>"))
				reader := bufio.NewReader(conn)
				line, err := reader.ReadString('\n')
				if err != nil {
					break
				}
				line = strings.TrimRight(line, "\r\n")
				fmt.Printf("nun:%+v", len(line))
				conn.Write([]byte(line + "\r\n"))
				conn.Write([]byte("aaa\r\n"))
				conn.Write([]byte("bbb\r\n"))
			}
		}()
	}
}

func main() {
	TestTelnet()
}
