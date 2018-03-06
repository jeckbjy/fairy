package test

import (
	"fairy"
	"fairy/util"
	"fairy/util/terminal"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

type Student struct {
	Account string `id:"id"`
}

func TestTable(t *testing.T) {
	data := "STRING\nAccount\nAccount\nJack\naaa"
	records := fairy.ReadTableFromString(data, &Student{})
	for _, cfg := range records.([]*Student) {
		fmt.Printf("%+v\n", cfg)
	}
}

func TestLog(t *testing.T) {
	fairy.Debug("%+v", "HelloWord!")
	fairy.Debug("%+v,%+v", "asdf", 1)
}

func TestSignal(t *testing.T) {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}

func TestTerminalColor(t *testing.T) {
	terminal.Foreground(terminal.Red)
	fmt.Sprintln("Red")
	terminal.Reset()
	terminal.Foreground(terminal.Blue)
	fmt.Sprintln("Blue")
	terminal.Reset()
}

func TestChannel(t *testing.T) {
	t.Log("start!!!")
	stopChan := make(chan bool)
	listener, err := net.Listen("tcp", ":8866")
	if err != nil {
		t.Log("aaaaaaa")
		return
	}
	go func() {
		select {
		case <-stopChan:
			t.Log("stop chan!")
			return
		default:
			t.Log("accept begin")
			_, err := listener.Accept()
			if err != nil {
				fmt.Println("accept fail!")
			}
			t.Log("accept end")
		}
	}()

	util.Sleep(1000)
	close(stopChan)
	util.Sleep(1000)
	listener.Close()
	util.Sleep(1000)
	close(stopChan)
	util.Sleep(1000)
	t.Error("end!!!")
}
