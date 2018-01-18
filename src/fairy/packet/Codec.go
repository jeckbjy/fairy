package packet

import (
	"encoding/binary"
	"fairy"
)

type Codec struct {
	buffer *fairy.Buffer
	flag   int
	data   []byte
	index  uint
}

func (self *Codec) CreateWriter(buffer *fairy.Buffer) {
	self.flag = 0
	// self.buffer = make([]byte, 1, 64)
}

func (self *Codec) CreateReader(buffer *fairy.Buffer) {
}

func (self *Codec) PushInt(value int) {
	self.PushInt64(int64(value))
}

func (self *Codec) PushInt64(value int64) {

}

func (self *Codec) PushUInt(value uint) {
	self.PushUInt64(uint64(value))
}

func (self *Codec) PushUInt64(value uint64) {
	self.index++
	if value != 0 {
		mask := 1 << self.index
		self.flag |= mask
	}
}

func (self *Codec) PushStr(value string) {
	self.index++
	if value != "" {
		mask := 1 << self.index
		self.flag |= mask
		// push string
	}
}

func (self *Codec) ReadUInt() uint {
	return uint(self.ReadUInt64())
}

func (self *Codec) ReadUInt64() uint64 {
	self.index++
	if self.HasFlag() {
		val, err := binary.ReadUvarint(self.buffer)
		if err != nil {
			panic(err)
		}
		return val
	}

	return 0
}

func (self *Codec) ReadStr() string {
	return ""
}

func (self *Codec) HasFlag() bool {
	return true
}

func (self *Codec) SetFlag() {

}
