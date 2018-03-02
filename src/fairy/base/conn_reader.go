package base

import "fairy"

type ConnReader struct {
	reader *fairy.Buffer
}

func (self *ConnReader) NewReader() {
	if self.reader == nil {
		self.reader = fairy.NewBuffer()
	}

	self.reader.Clear()
}

func (self *ConnReader) Append(data []byte) {
	self.reader.Append(data)
}

func (self *ConnReader) Reader() *fairy.Buffer {
	return self.reader
}
