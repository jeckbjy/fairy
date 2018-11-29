package exit

import (
	"os"
	"os/signal"
	"syscall"
)

// Callback listener for exit
type Callback func()

var handlers []Callback

// Add 注册回调
func Add(cb Callback) {
	handlers = append(handlers, cb)
}

// Wait wait for graceful exit
func Wait(sig ...os.Signal) {
	exitSig := make(chan os.Signal)
	if len(sig) > 0 {
		signal.Notify(exitSig, sig...)
	} else {
		signal.Notify(exitSig, syscall.SIGINT, syscall.SIGTERM)
	}
	<-exitSig
	for _, cb := range handlers {
		cb()
	}
}
