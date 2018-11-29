package filters

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
type TelnetCB func(conn fairy.IConn, str string)

// 默认回调
func defaultCB(conn fairy.IConn, str string) {
	handler := fairy.GetDispatcher().GetHandlerByName(TelnetMsgKey)
	if handler != nil {
		pkt := packet.NewBase()
		pkt.SetName(TelnetMsgKey)
		pkt.SetMessage(str)
		fairy.GetExecutor().Dispatch(0, func() {
			fairy.InvokeHandler(conn, pkt, handler)
		})
	} else {
		log.Error("cannot find telnet handler!")
	}
}

// NewTelnet 创建TelnetFilter
func NewTelnet() *TelnetFilter {
	return &TelnetFilter{prompt: "fairy>", cb: defaultCB}
}

// TelnetFilter 远程连接
type TelnetFilter struct {
	base.Filter
	prompt string
	cb     TelnetCB
}

// SetPrompt 设置提示符
func (tf *TelnetFilter) SetPrompt(value string) {
	tf.prompt = value
}

// SetCallback 设置回调函数
func (tf *TelnetFilter) SetCallback(cb TelnetCB) {
	tf.cb = cb
}

func (tf *TelnetFilter) Name() string {
	return "TelnetFilter"
}

func (self *TelnetFilter) HandleOpen(ctx fairy.IFilterCtx) {
	if self.prompt != "" {
		conn := ctx.GetConn()
		buffer := fairy.NewBuffer()
		buffer.Append([]byte(self.prompt))
		conn.Write(buffer)
	}

	ctx.Next()
}

func (self *TelnetFilter) HandleRead(ctx fairy.IFilterCtx) {
	// parse \r\n
	if buffer, ok := ctx.GetData().(*fairy.Buffer); ok {
		result, err := buffer.ReadLine()
		if err == nil {
			str := result.String()
			ctx.SetData(str)
			// 默认行为
			if self.cb != nil {
				self.cb(ctx.GetConn(), str)
			}
		}
	}

	ctx.Next()
}

func (tf *TelnetFilter) HandleWrite(ctx fairy.IFilterCtx) {
	// write message
	if str, ok := ctx.GetData().(string); ok {
		buffer := fairy.NewBuffer()
		buffer.Append([]byte(str))
		if !strings.HasSuffix(str, "\r\n") {
			buffer.Append([]byte("\r\n"))
		}

		if tf.prompt != "" {
			buffer.Append([]byte(tf.prompt))
		}

		ctx.SetData(buffer)
	}

	ctx.Next()
}
