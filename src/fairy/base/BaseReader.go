package base

import "fairy"

type BaseReader struct {
	reader *fairy.Buffer
}

func (self *BaseReader) NewReader() {
	if self.reader == nil {
		self.reader = fairy.NewBuffer()
	}

	self.reader.Clear()
}

func (self *BaseReader) Append(data []byte) {
	self.reader.Append(data)
}

func (self *BaseReader) Reader() *fairy.Buffer {
	return self.reader
}
