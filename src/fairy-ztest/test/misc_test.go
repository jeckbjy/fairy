package test

import (
	"fairy"
	"fairy/util"
	"fairy/util/terminal"
	"fmt"
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

func TestTimer(t *testing.T) {
	tt := util.Now()
	fairy.StartTimer(10, func(timer *fairy.Timer) {
		fairy.Debug("OnTimer out:%+v", util.Now()-tt)
	})

	util.Sleep(10000)
}
