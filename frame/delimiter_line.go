package frame

import (
	"github.com/jeckbjy/fairy"
)

func NewLine() fairy.Frame {
	return NewLineEx("\n")
}

func NewLineEx(delimiter string) fairy.Frame {
	frame := &LineFrame{}
	frame.delimiter = delimiter
	return frame
}

// \r\n or \n
type LineFrame struct {
	delimiter string
}

func (self *LineFrame) Encode(buffer *fairy.Buffer) error {
	buffer.Append([]byte(self.delimiter))
	return nil
}

func (self *LineFrame) Decode(buffer *fairy.Buffer) (*fairy.Buffer, error) {
	return buffer.ReadLine()
}
