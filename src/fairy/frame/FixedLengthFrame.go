package frame

import (
	"errors"
	"fairy"
	"fairy/util"
	"io"
	"math"
)

func NewFixedLengthFrame(headLen uint, littleEndian bool) *FixedLengthFrame {
	frame := &FixedLengthFrame{}
	frame.headLen = headLen
	frame.littleEndian = littleEndian
	frame.CalcMaxLen()
	return frame
}

type FixedLengthFrame struct {
	headLen      uint
	maxLen       uint
	littleEndian bool
}

func (self *FixedLengthFrame) CalcMaxLen() {
	switch self.headLen {
	case 1:
		self.maxLen = math.MaxUint8
	case 2:
		self.maxLen = math.MaxUint16
	case 4:
		self.maxLen = math.MaxUint32
	default:
		self.headLen = 4
		self.maxLen = math.MaxUint32
	}
}

func (self *FixedLengthFrame) Encode(buffer *fairy.Buffer) error {
	leng := buffer.Length()
	if uint(leng) > self.maxLen {
		return errors.New("buffer overflow in FixedLengthFrame!")
	}

	data := make([]byte, self.headLen)
	switch self.headLen {
	case 1:
		data[0] = byte(leng)
	case 2:
		util.PutUint16(data, uint16(leng), self.littleEndian)
	case 4:
		util.PutUint32(data, uint32(leng), self.littleEndian)
	}
	// 前边插入
	buffer.Prepend(data)
	return nil
}

func (self *FixedLengthFrame) Decode(buffer *fairy.Buffer) (*fairy.Buffer, error) {
	// read data length
	data := make([]byte, self.headLen)
	if _, err := buffer.Read(data); err != nil {
		return nil, errors.New("FixedLengthFrame read lenth fail!")
	}

	var count uint
	// check len and read
	switch self.headLen {
	case 1:
		count = uint(data[0])
	case 2:
		count = uint(util.GetUint16(data, self.littleEndian))
	case 4:
		count = uint(util.GetUint32(data, self.littleEndian))
	}

	// check has enough data
	if !buffer.HasRemain(int(count)) {
		return nil, errors.New("FixedLengthFrame read data fail!")
	}

	result := fairy.NewBuffer()

	// discard length
	buffer.Discard()
	buffer.Seek(int(count), io.SeekStart)
	buffer.Split(result)
	return result, nil
}
