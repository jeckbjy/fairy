package packet

import (
	"encoding/binary"

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

// Reader 用于通过bit位解析消息头
type Reader struct {
	buffer *fairy.Buffer // 缓存
	flag   uint          // 读取的标识
	mask   uint          // 将要检测的掩码
}

func (reader *Reader) next() bool {
	reader.mask <<= 1
	return (reader.flag & reader.mask) != 0
}

func (reader *Reader) Read(data []byte) (int, error) {
	return reader.buffer.Read(data)
}

func (reader *Reader) GetStr() (string, error) {
	if !reader.next() {
		return "", nil
	}
	// 先读取长度
	len, err := binary.ReadUvarint(reader.buffer)
	if err != nil {
		return "", err
	}

	// 读取字符串
	str := make([]byte, len)
	_, serr := reader.buffer.Read(str)
	if serr != nil {
		return "", err
	}

	return string(str), nil
}

func (reader *Reader) GetVar64() (int64, error) {
	if !reader.next() {
		return 0, nil
	}

	return binary.ReadVarint(reader.buffer)
}

func (reader *Reader) GetUVar64() (uint64, error) {
	if !reader.next() {
		return 0, nil
	}

	return binary.ReadUvarint(reader.buffer)
}

func (reader *Reader) GetUInt16LE() (uint16, error) {
	if !reader.next() {
		return 0, nil
	}

	data := make([]byte, 2)
	if _, err := reader.Read(data); err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint16(data)
	return value, nil
}

func (reader *Reader) GetUInt16BE() (uint16, error) {
	if !reader.next() {
		return 0, nil
	}

	data := make([]byte, 2)
	if _, err := reader.Read(data); err != nil {
		return 0, err
	}

	value := binary.BigEndian.Uint16(data)
	return value, nil
}

func (reader *Reader) GetUInt32LE() (uint32, error) {
	if !reader.next() {
		return 0, nil
	}

	data := make([]byte, 4)
	if _, err := reader.Read(data); err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint32(data)
	return value, nil
}

func (reader *Reader) GetUInt32BE() (uint32, error) {
	if !reader.next() {
		return 0, nil
	}

	data := make([]byte, 4)
	if _, err := reader.Read(data); err != nil {
		return 0, err
	}
	value := binary.BigEndian.Uint32(data)
	return value, nil
}

func (reader *Reader) GetUInt64LE() (uint64, error) {
	if !reader.next() {
		return 0, nil
	}

	data := make([]byte, 8)
	if _, err := reader.Read(data); err != nil {
		return 0, err
	}
	value := binary.LittleEndian.Uint64(data)
	return value, nil
}

func (reader *Reader) GetUInt64BE() (uint64, error) {
	if !reader.next() {
		return 0, nil
	}

	data := make([]byte, 4)
	if _, err := reader.Read(data); err != nil {
		return 0, err
	}
	value := binary.BigEndian.Uint64(data)
	return value, nil
}

func (reader *Reader) GetInt16LE() (int16, error) {
	v, e := reader.GetUInt16LE()
	return int16(v), e
}

func (reader *Reader) GetInt16BE() (int16, error) {
	v, e := reader.GetUInt16BE()
	return int16(v), e
}

func (reader *Reader) GetInt32LE() (int32, error) {
	v, e := reader.GetUInt32LE()
	return int32(v), e
}

func (reader *Reader) GetInt32BE() (int32, error) {
	v, e := reader.GetUInt32BE()
	return int32(v), e
}

func (reader *Reader) GetInt64LE() (int64, error) {
	v, e := reader.GetUInt64LE()
	return int64(v), e
}

func (reader *Reader) GetInt64BE() (int64, error) {
	v, e := reader.GetUInt64BE()
	return int64(v), e
}
