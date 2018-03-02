package frame

import (
	"fairy"
)

func NewLineFrame() fairy.Frame {
	return NewLineFrameEx("\n")
}

func NewLineFrameEx(delimiter string) fairy.Frame {
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
	// delimiterCount := 1
	// pos := buffer.IndexOf("\n")
	// if pos == -1 {
	// 	return nil, errors.New("donnot find delimiter!")
	// }

	// // check \r
	// if pos > 0 {
	// 	buffer.Seek(-1, io.SeekCurrent)
	// 	if ch, _ := buffer.ReadByte(); ch == '\r' {
	// 		delimiterCount = 2
	// 		pos -= 1
	// 		buffer.Seek(-1, io.SeekCurrent)
	// 	}
	// }

	// result := fairy.NewBuffer()
	// // read data
	// buffer.Seek(pos, io.SeekStart)
	// buffer.Split(result)

	// // remove demilter
	// buffer.Seek(delimiterCount, io.SeekStart)
	// buffer.Discard()

	// fairy.Debug("msgs:%+v", result.String())
	// fairy.Debug("left:%+v", buffer.String())

	// return result, nil
}
