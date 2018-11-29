package fairy

import (
	"fmt"
	"math"
	"reflect"
)

var gDispatcher = NewDispatcher()

func GetDispatcher() *Dispatcher {
	return gDispatcher
}

func NewDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{}
	dispatcher.nameMap = make(map[string]*Handler)
	return dispatcher
}

// RegisterHandler 注册消息回调
func RegisterHandler(key interface{}, cb HandlerCB) *Handler {
	return GetDispatcher().Register(key, cb, QueueMainID)
}

// InvokeHandler 调用Handler
func InvokeHandler(conn IConn, packet IPacket, handler *Handler) {
	hctx := HandlerCtx{}
	hctx.Init(conn, packet, handler, GetDispatcher().Middlewares())
	hctx.Next()
}

// HandlerCtx 支持中间件,在Dispatcher中注册
type HandlerCtx struct {
	IConn
	packet  IPacket
	args    []interface{} // 附加参数
	handler *Handler      // 回调
	chain   []HandlerCB   // 中间件
	index   int           // 索引
}

func (ctx *HandlerCtx) Init(conn IConn, packet IPacket, handler *Handler, chain []HandlerCB) {
	ctx.IConn = conn
	ctx.packet = packet
	ctx.handler = handler
	ctx.chain = chain
	ctx.index = -1
}

func (ctx *HandlerCtx) Packet() IPacket {
	return ctx.packet
}

func (ctx *HandlerCtx) MsgID() uint {
	return ctx.packet.GetId()
}

func (ctx *HandlerCtx) MsgName() string {
	return ctx.packet.GetName()
}

func (ctx *HandlerCtx) Message() interface{} {
	return ctx.packet.GetMessage()
}

func (ctx *HandlerCtx) Len() int {
	return len(ctx.args)
}

func (ctx *HandlerCtx) Add(v interface{}) {
	ctx.args = append(ctx.args, v)
}

func (ctx *HandlerCtx) Get(index int) interface{} {
	return ctx.args[index]
}

func (ctx *HandlerCtx) Next() {
	ctx.index++

	if ctx.index < len(ctx.chain) {
		cb := ctx.chain[ctx.index]
		cb(ctx)
	} else if ctx.index == len(ctx.chain) {
		ctx.handler.Func(ctx)
		// finish
		ctx.index++
	}
}

// HandlerCB 消息回调
type HandlerCB func(ctx *HandlerCtx)

// Handler 消息处理
type Handler struct {
	Func    HandlerCB
	QueueId uint
	Data    interface{}
}

// Dispatcher 消息回调管理,id不能超过uint16
// 支持添加中间件
type Dispatcher struct {
	nameMap     map[string]*Handler
	idArray     []*Handler
	middlewares []HandlerCB
}

func (d *Dispatcher) Use(cb HandlerCB) {
	d.middlewares = append(d.middlewares, cb)
}

func (d *Dispatcher) Middlewares() []HandlerCB {
	return d.middlewares
}

// Register 注册消息回调,可以指定队列
func (d *Dispatcher) Register(key interface{}, cb HandlerCB, queueId uint) *Handler {
	id := uint(0)
	name := ""
	switch key.(type) {
	case int:
		id = uint(key.(int))
	case uint:
		id = key.(uint)
	case string:
		name = key.(string)
	default:
		// must be struct!!!
		// example:Register(&LoginReq{}, handler) or Register(LoginReq{}, handler)
		rtype := reflect.TypeOf(key)
		if rtype.Kind() == reflect.Ptr {
			rtype = rtype.Elem()
		}

		if rtype.Kind() == reflect.Struct {
			name = rtype.Name()
		} else {
			panic(fmt.Sprintf("register handler fail,bad key type!key=%+v", key))
		}
	}

	handler := &Handler{Func: cb, QueueId: queueId}

	if id > 0 {
		// 通过id注册
		if id > math.MaxUint16 {
			panic("handler id overflow")
		}

		if int(id) >= len(d.idArray) {
			targets := make([]*Handler, id+1)
			copy(targets, d.idArray)
			d.idArray = targets
		}
		d.idArray[id] = handler
	} else if name != "" {
		// 通过名字注册
		d.nameMap[name] = handler
	}

	return handler
}

func (d *Dispatcher) GetHandler(id uint, name string) *Handler {
	if id > 0 {
		// 通过id查找
		if int(id) < len(d.idArray) {
			return d.idArray[id]
		}

		return nil
	}

	return d.nameMap[name]
}

func (d *Dispatcher) GetHandlerById(id uint) *Handler {
	if int(id) < len(d.idArray) {
		return d.idArray[id]
	}

	return nil
}

func (d *Dispatcher) GetHandlerByName(name string) *Handler {
	return d.nameMap[name]
}
