package packet

import (
	"encoding/binary"
	"fmt"

	"github.com/jeckbjy/fairy"
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
