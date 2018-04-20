package packet

import (
	"encoding/binary"

	"github.com/jeckbjy/fairy"
)

func NewWriter(buffer *fairy.Buffer) *Writer {
	writer := &Writer{buffer: buffer}
	writer.data = make([]byte, 1, 255)
	writer.vbuf = make([]byte, 10)
	return writer
}

type Writer struct {
	buffer *fairy.Buffer
	flag   uint
	mask   uint
	data   []byte
	vbuf   []byte
}

func (self *Writer) Flush() {
	self.data[0] = byte(self.flag)
	self.buffer.Append(self.data)
}

func (self *Writer) putUvarint(v uint64) {
	n := binary.PutUvarint(self.vbuf, v)
	self.data = append(self.data, self.vbuf[:n]...)
}

func (self *Writer) PutId(v uint) {
	// uint16 and no mask
	binary.LittleEndian.PutUint16(self.vbuf, uint16(v))
	self.data = append(self.data, self.vbuf[:2]...)
}

func (self *Writer) PutUint(v uint) {
	self.PutUint64(uint64(v))
}

func (self *Writer) PutUint64(v uint64) {
	self.mask <<= 1
	if v != 0 {
		self.flag |= self.mask
		self.putUvarint(v)
	}
}

func (self *Writer) PutStr(v string) {
	self.mask <<= 1
	if v != "" {
		self.flag |= self.mask
		self.putUvarint(uint64(len(v)))
		self.data = append(self.data, []byte(v)...)
	}
}
