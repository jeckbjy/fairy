package filter

import (
	"fairy"
	"fairy/base"
	"strings"
)

func NewTelnet() *TelnetFilter {
	return NewTelnetFilterEx("fairy>")
}

func NewTelnetFilterEx(prompt string) *TelnetFilter {
	f := &TelnetFilter{}
	f.Prompt = prompt
	return f
}

type TelnetFilter struct {
	base.BaseFilter
	Prompt string
}

func (self *TelnetFilter) HandleOpen(ctx fairy.FilterContext) fairy.FilterAction {
	if self.Prompt != "" {
		conn := ctx.GetConnection()
		buffer := fairy.NewBuffer()
		buffer.Append([]byte(self.Prompt))
		conn.Write(buffer)
	}

	return ctx.GetNextAction()
}

func (self *TelnetFilter) HandleRead(ctx fairy.FilterContext) fairy.FilterAction {
	// parse \r\n
	if buffer, ok := ctx.GetMessage().(*fairy.Buffer); ok {
		result, err := buffer.ReadLine()
		if err == nil {
			str := result.String()
			ctx.SetMessage(str)
			// fmt.Printf("%+v", str)
			// conn := ctx.GetConnection()
			// conn.Send(str)
		}
	}

	return ctx.GetNextAction()
}

func (self *TelnetFilter) HandleWrite(ctx fairy.FilterContext) fairy.FilterAction {
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
