package frame

import (
	"errors"
	"fairy"
	"io"
)

func NewLineFrame() fairy.Frame {
	return NewDelimiterFrame("\n")
}

func NewDelimiterFrame(delimiter string) fairy.Frame {
	frame := &DelimiterFrame{}
	frame.delimiter = delimiter
	if delimiter == "" {
		panic("delimiter is null string!")
	}
	return frame
}

type DelimiterFrame struct {
	delimiter string
}

func (self *DelimiterFrame) Encode(buffer *fairy.Buffer) error {
	// change to write??
	buffer.Append([]byte(self.delimiter))
	return nil
}

func (self *DelimiterFrame) Decode(buffer *fairy.Buffer) (*fairy.Buffer, error) {
	pos := buffer.IndexOf(self.delimiter)
	if pos == -1 {
		return nil, errors.New("donnot find delimiter!")
	}

	result := fairy.NewBuffer()
	// seek data, if pos == 0 ??
	buffer.Seek(pos, io.SeekCurrent)
	buffer.Split(result)

	// seek delimiter
	buffer.Seek(len(self.delimiter), io.SeekCurrent)
	buffer.Discard()

	return result, nil
}
