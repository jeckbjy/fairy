package fairy

import "github.com/jeckbjy/fairy/util"

func WaitExit() {
	util.GetExit().Wait()
}

// RegisterExit 注册退出
func RegisterExit(hander util.ExitHandler) {
	util.GetExit().Register(hander)
}

// RegisterMessage 注册消息
func RegisterMessage(msg interface{}, args ...interface{}) {
	GetRegistry().Register(msg, args...)
}

// RegisterHandler 注册简单的消息回调函数
func RegisterHandler(key interface{}, cb HandlerCB) {
	RegisterHandlerEx(key, cb, 0)
}

// RegisterHandlerEx 注册消息回调函数
func RegisterHandlerEx(key interface{}, cb HandlerCB, queueID int) {
	GetDispatcher().Regsiter(key, &HandlerHolder{cb: cb, queueID: queueID})
}

// RegisterMiddleware 注册中间件
func RegisterMiddleware(cb HandlerCB) {
	GetDispatcher().Use(cb)
}
