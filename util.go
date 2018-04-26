package fairy

import "github.com/jeckbjy/fairy/util"

func WaitExit() {
	util.GetExit().Wait()
}

func RegisterExit(hander util.ExitHandler) {
	util.GetExit().Register(hander)
}

// 注册消息
func RegisterMessage(msg interface{}, args ...interface{}) {
	GetRegistry().Register(msg, args...)
}

//////////////////////////////////////////////////////////
// 注册回调函数
//////////////////////////////////////////////////////////
func RegisterHandler(key interface{}, cb HandlerCB) {
	RegisterHandlerEx(key, cb, 0)
}

func RegisterHandlerEx(key interface{}, cb HandlerCB, queueId int) {
	GetDispatcher().Regsiter(key, &HandlerHolder{cb: cb, queueID: queueId})
}

func RegisterUncaughtHandler(cb HandlerCB) {
	GetDispatcher().SetUncaughtHandler(&HandlerHolder{cb: cb, queueID: 0})
}
