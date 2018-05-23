package filter

import (
	"strings"

	"github.com/jeckbjy/fairy"
	"github.com/jeckbjy/fairy/base"
	"github.com/jeckbjy/fairy/log"
	"github.com/jeckbjy/fairy/packet"
)

const (
	// TelnetMsgKey 消息唯一标识
	TelnetMsgKey = "TelnetMsg"
)

// TelnetCB 读处理回调
type TelnetCB func(conn fairy.Conn, str string)

// 默认的出来回调
func defaultCB(conn fairy.Conn, str string) {
	handler := fairy.GetDispatcher().GetHandlerByName(TelnetMsgKey)
	if handler != nil {
		pkt := packet.NewBase()
		pkt.SetName(TelnetMsgKey)
		pkt.SetMessage(str)
		ctx := fairy.NewHandlerCtx(conn, pkt, handler, fairy.GetDispatcher().Middlewares())
		fairy.GetExecutor().Dispatch(ctx)
	} else {
		log.Error("cannot find telnet handler!")
	}
}

// NewTelnet 创建TelnetFilter
func NewTelnet() fairy.Filter {
	return NewTelnetEx(defaultCB, "fairy>")
}

// NewTelnetEx 带有参数的创建TelnetFilter
func NewTelnetEx(cb TelnetCB, prompt string) fairy.Filter {
	f := &zTelnetFilter{}
	f.Prompt = prompt
	f.cb = cb
	return f
}

/**
 * zTelnetFilter 使用方法
 * func telnet_cb(conn fairy.Conn, pkt fairy.Packet) {
 * 		str := pkt.GetMessage().(str)
 * }
 *
 * 1:默认通过注册回调函数,实现调用,默认会在主线程中处理
 * fairy.RegisterHandler(filter.TelnetMsgKey, telnet_cb)
 * 2:创建时注册回调函数
 */
type zTelnetFilter struct {
	base.BaseFilter
	Prompt string
	cb     TelnetCB
}

func (self *zTelnetFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	if self.Prompt != "" {
		conn := ctx.GetConn()
		buffer := fairy.NewBuffer()
		buffer.Append([]byte(self.Prompt))
		conn.Write(buffer)
	}

	return ctx.GetNextAction()
}

func (self *zTelnetFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	// parse \r\n
	if buffer, ok := ctx.GetMessage().(*fairy.Buffer); ok {
		result, err := buffer.ReadLine()
		if err == nil {
			str := result.String()
			ctx.SetMessage(str)
			// 默认行为
			if self.cb != nil {
				self.cb(ctx.GetConn(), str)
			}
		}
	}

	return ctx.GetNextAction()
}

func (self *zTelnetFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
	// write message
	if str, ok := ctx.GetMessage().(string); ok {
		buffer := fairy.NewBuffer()
		buffer.Append([]byte(str))
		if !strings.HasSuffix(str, "\r\n") {
			buffer.Append([]byte("\r\n"))
		}
		if self.Prompt != "" {
			buffer.Append([]byte(self.Prompt))
		}

		ctx.SetMessage(buffer)
	}

	return ctx.GetNextAction()
}
