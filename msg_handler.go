package fairy

type HandlerCtx struct {
	Conn
	Packet
	args []interface{}
}

func (ctx *HandlerCtx) Set(args ...interface{}) {
	ctx.args = append(ctx.args, args...)
}

func (ctx *HandlerCtx) Get(idx int) interface{} {
	return ctx.args[idx]
}

type HandlerCB func(ctx *HandlerCtx)

// Handler 消息回调接口
type Handler interface {
	QueueID() int
	Invoke(ctx *HandlerCtx)
}

// HandlerHolder 持有callback
type HandlerHolder struct {
	queueID int
	cb      HandlerCB
}

// QueueID 返回queueid
func (holder *HandlerHolder) QueueID() int {
	return holder.queueID
}

// Invoke 调用callback
func (holder *HandlerHolder) Invoke(ctx *HandlerCtx) {
	holder.cb(ctx)
}
