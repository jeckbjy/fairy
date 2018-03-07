package frame

import (
	"errors"
	"fairy"
	"io"
)

func NewDelimiter(delimiter string) fairy.Frame {
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
	buffer.Append([]byte(self.delimiter))
	return nil
}

func (self *DelimiterFrame) Decode(buffer *fairy.Buffer) (*fairy.Buffer, error) {
	pos := buffer.IndexOf(self.delimiter)
	if pos == -1 {
		return nil, errors.New("donnot find delimiter!")
	}

	result := fairy.NewBuffer()
	// read data
	buffer.Seek(pos, io.SeekStart)
	buffer.Split(result)

	// remove delimiter
	buffer.Seek(len(self.delimiter), io.SeekStart)
	buffer.Discard()

	return result, nil
}
