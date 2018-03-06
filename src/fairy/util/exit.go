package util

import (
	"os"
)

func RegisterExit(hander ExitHandler) {
	GetExit().Register(hander)
}

func WaitExit() {
	GetExit().Wait()
}

var gExit *Exit

func GetExit() *Exit {
	if gExit == nil {
		gExit = &Exit{}
		gExit.Create()
	}

	return gExit
}

type ExitHandler interface {
	OnExit()
}

type Exit struct {
	sig      chan os.Signal
	handlers []ExitHandler
}

func (self *Exit) Create() {
	self.sig = make(chan os.Signal)
}

func (self *Exit) Register(handler ExitHandler) {
	self.handlers = append(self.handlers, handler)
}

func (self *Exit) Run() {
	go self.Wait()
}

func (self *Exit) Wait() {
	<-self.sig
	for _, handler := range self.handlers {
		handler.OnExit()
	}
}
