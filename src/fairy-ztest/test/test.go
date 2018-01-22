package test

import (
	"fairy"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Student struct {
	Account string `id:"id"`
}

func TestTable() {
	data := "STRING\nAccount\nAccount\nJack\naaa"
	records := fairy.ReadTableFromString(data, &Student{})
	for _, cfg := range records.([]*Student) {
		fmt.Printf("%+v\n", cfg)
	}
}

func TestLog() {
	fairy.Debug("%+v", "HelloWord!")
	fairy.Debug("%+v,%+v", "asdf", 1)
}

func TestSignal() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-c
}
