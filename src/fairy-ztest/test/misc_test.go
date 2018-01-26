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
	var gTimerStart = util.Now()
	fairy.StartTimer(util.FromSec(2), func(timer *fairy.Timer) {
		diff := util.Now() - gTimerStart
		if diff/1000 != 2 {
			t.Error("timer fail!")
		} else {
			t.Log("timer succeed!")
		}

		os.Exit(0)
	})

	fairy.WaitExit()
}
