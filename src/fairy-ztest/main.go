package main

import (
	"fairy"
	"fairy/util"
	"os"
)

func TestTimer() {
	var gTimerStart = util.Now()
	fairy.StartTimer(util.FromSec(20), func(timer *fairy.Timer) {
		diff := util.Now() - gTimerStart
		if diff/1000 != 20 {
			fairy.Error("timer fail!")
		} else {
			fairy.Debug("timer succeed!")
		}

		os.Exit(0)
	})

	fairy.WaitExit()
}

func main() {
	TestTimer()
	// mode := flag.String("m", "server", "server mode")
	// flag.Parse()
	// if *mode == "server" {
	// 	echo.StartServer()
	// } else {
	// 	echo.StartClient()
	// }
}
