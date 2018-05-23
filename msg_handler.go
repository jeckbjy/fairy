package fairy

// NewHandlerCtx 创建Ctx
func NewHandlerCtx(conn Conn, pkt Packet, hanlder Handler, middlewares HandlerChain) *HandlerCtx {
	return &HandlerCtx{Conn: conn, Packet: pkt, handler: hanlder, middlewares: middlewares, index: -1}
}

// HandlerCtx 用于Handler传递信息
// TODO:need pool?
type HandlerCtx struct {
	Conn
	Packet
	args        []interface{}
	handler     Handler
	middlewares HandlerChain // 中间件
	index       int8         // 初始-1
}

// PushArgs 添加参数
func (ctx *HandlerCtx) PushArgs(args ...interface{}) {
	ctx.args = append(ctx.args, args...)
}

// Args 获取参数
func (ctx *HandlerCtx) Args(idx int) interface{} {
	return ctx.args[idx]
}

// Reset 重置数据
func (ctx *HandlerCtx) Reset() {
	ctx.Conn = nil
	ctx.Packet = nil
	ctx.args = nil
	ctx.index = -1
}

// 返回队列ID
func (ctx *HandlerCtx) QueueID() int {
	if ctx.handler != nil {
		return ctx.handler.QueueID()
	}

	return 0
}

// Next 执行下一个
func (ctx *HandlerCtx) Next() {
	ctx.index++

	if ctx.index < int8(len(ctx.middlewares)) {
		for s := int8(len(ctx.middlewares)); ctx.index < s; ctx.index++ {
			ctx.middlewares[ctx.index](ctx)
		}
	} else {
		if ctx.handler != nil {
			ctx.handler.Invoke(ctx)
		}
		ctx.index++
	}
}

// Process 实现Event接口
func (ctx *HandlerCtx) Process() {
	ctx.Next()
}

type HandlerCB func(ctx *HandlerCtx)
type HandlerChain []HandlerCB

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
