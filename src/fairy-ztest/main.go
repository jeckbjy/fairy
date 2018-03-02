package main

import (
	"bufio"
	"fairy-ztest/echo"
	"flag"
	"fmt"
	"net"
	"strings"
)

func TestFairy() {
	pside := flag.String("s", "server", "server or client side")
	pnetmode := flag.String("n", "ws", "network mode,tcp,ws,kcp")
	pmsgmode := flag.String("m", "json", "proto mode,json,protobuf,sproto,bson")
	flag.Parse()

	if *pside == "server" {
		echo.StartServer(*pnetmode, *pmsgmode)
	} else {
		echo.StartClient(*pnetmode, *pmsgmode)
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
	TestFairy()
	// TestTelnet()
}
