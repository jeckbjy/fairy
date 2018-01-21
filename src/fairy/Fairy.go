package fairy

import (
	"fairy/util"
)

func RegisterExit(hander util.ExitHandler) {
	util.GetExit().Register(hander)
}

func WaitExit() {
	util.GetExit().Wait()
}
