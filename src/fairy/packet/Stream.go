package packet

import (
	"encoding/binary"
	"fairy"
	"fmt"
)

func NewReader(buffer *fairy.Buffer) *Reader {
	flag, err := buffer.ReadByte()
	if err != nil {
		return nil
	}

	reader := &Reader{buffer: buffer, flag: uint(flag)}
	return reader
}

type Reader struct {
	buffer *fairy.Buffer
	flag   uint
	mask   uint
}

func (self *Reader) next() bool {
	self.mask <<= 1
	return (self.flag & self.mask) != 0
}

func (self *Reader) GetId() (uint, error) {
	buf := make([]byte, 2)
	if _, err := self.buffer.Read(buf); err != nil {
		return 0, err
	}

	ret := binary.LittleEndian.Uint16(buf)
	if ret == 0 {
		return 0, fmt.Errorf("id cannot zero!")
	}

	return uint(ret), nil
}

func (self *Reader) GetUint() (uint, error) {
	if self.next() {
		ret, err := binary.ReadUvarint(self.buffer)
		return uint(ret), err
	}

	return 0, nil
}

func (self *Reader) GetUint64() (uint64, error) {
	if self.next() {
		ret, err := binary.ReadUvarint(self.buffer)
		return ret, err
	}

	return 0, nil
}

func (self *Reader) GetStr() (string, error) {
	if self.next() {
		len, err := binary.ReadUvarint(self.buffer)
		if err != nil {
			return "", err
		}

		str := make([]byte, len)
		_, err = self.buffer.Read(str)
		if err != nil {
			return "", err
		}

		return string(str), nil
	}

	return "", nil
}

/////////////////////////////////////////////////////////
// writer
/////////////////////////////////////////////////////////

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
