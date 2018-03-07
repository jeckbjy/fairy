package test

import (
	"bufio"
	"fairy"
	"fairy/filter"
	"fairy/tcp"
	"fmt"
	"net"
	"strings"
	"testing"
)

func TelnetDemo() {
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

func TestTelnet(t *testing.T) {
	tran := tcp.NewTransport()
	tran.AddFilters(
		filter.NewTelnet())
	tran.Listen(":8080", 0)
	fairy.WaitExit()
}
