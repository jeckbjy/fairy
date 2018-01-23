package test

import (
	"fairy"
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
