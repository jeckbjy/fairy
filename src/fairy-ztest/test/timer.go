package test

import (
	"fairy/log"
	"fairy/timer"
	"fairy/util"
)

func TestTimer() {
	tt := util.Now()
	timer.Start(10000, func(t *timer.Timer) {
		log.Debug("OnTimer out:%+v", util.Now()-tt)
	})

	util.Sleep(100000)
}
