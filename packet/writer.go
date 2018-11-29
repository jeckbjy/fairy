package packet

import (
	"encoding/binary"

	"github.com/jeckbjy/fairy"
)

// NewWriter 创建writer
func NewWriter(buffer *fairy.Buffer) *Writer {
	writer := &Writer{buffer: buffer}
	writer.data = make([]byte, 1, 255)
	writer.vbuf = make([]byte, 10)
	return writer
}

// Writer 用于序列化
type Writer struct {
	buffer *fairy.Buffer
	flag   uint   // 最终标识
	mask   uint   // 当前掩码
	data   []byte // 数据缓存
	vbuf   []byte // PutUvarint 缓存
}

// Flush 将缓存中数据写入buffer
func (writer *Writer) Flush() {
	writer.data[0] = byte(writer.flag)
	writer.buffer.Append(writer.data)
}

func (writer *Writer) next() {
	writer.mask <<= 1
}

func (writer *Writer) use() {
	writer.flag |= writer.mask
}

func (writer *Writer) appendVar(x int64) {
	n := binary.PutVarint(writer.vbuf, x)
	writer.Append(writer.vbuf[:n])
}

func (writer *Writer) appendUVar(x uint64) {
	n := binary.PutUvarint(writer.vbuf, x)
	writer.Append(writer.vbuf[:n])
}

// Append 追加但不设置mask
func (writer *Writer) Append(b []byte) {
	writer.data = append(writer.data, b...)
}

func (writer *Writer) PutStr(v string) {
	writer.next()
	if v != "" {
		writer.use()
		writer.appendUVar(uint64(len(v)))
		writer.Append([]byte(v))
	}
}

func (writer *Writer) PutUVar64(v uint64) {
	writer.next()
	if v != 0 {
		writer.use()
		writer.appendUVar(v)
	}
}

func (writer *Writer) PutVar64(v int64) {
	writer.next()
	if v != 0 {
		writer.use()
		writer.appendVar(v)
	}
}

func (writer *Writer) PutUInt16LE(v uint16) {
	writer.next()
	if v != 0 {
		writer.use()
		binary.LittleEndian.PutUint16(writer.vbuf, v)
		writer.Append(writer.vbuf[:2])
	}
}

func (writer *Writer) PutUInt16BE(v uint16) {
	writer.next()
	if v != 0 {
		writer.use()
		binary.BigEndian.PutUint16(writer.vbuf, v)
		writer.Append(writer.vbuf[:2])
	}
}

func (writer *Writer) PutUInt32LE(v uint32) {
	writer.next()
	if v != 0 {
		writer.use()
		binary.LittleEndian.PutUint32(writer.vbuf, v)
		writer.Append(writer.vbuf[:4])
	}
}

func (writer *Writer) PutUInt32BE(v uint32) {
	writer.next()
	if v != 0 {
		writer.use()
		binary.BigEndian.PutUint32(writer.vbuf, v)
		writer.Append(writer.vbuf[:4])
	}
}

func (writer *Writer) PutUInt64LE(v uint64) {
	writer.next()
	if v != 0 {
		writer.use()
		binary.LittleEndian.PutUint64(writer.vbuf, v)
		writer.Append(writer.vbuf[:8])
	}
}

func (writer *Writer) PutUInt64BE(v uint64) {
	writer.next()
	if v != 0 {
		writer.use()
		binary.BigEndian.PutUint64(writer.vbuf, v)
		writer.Append(writer.vbuf[:8])
	}
}

func (writer *Writer) PutInt16LE(v int16) {
	writer.PutUInt16LE(uint16(v))
}

func (writer *Writer) PutInt16BE(v int16) {
	writer.PutUInt16BE(uint16(v))
}

func (writer *Writer) PutInt32LE(v int32) {
	writer.PutUInt32LE(uint32(v))
}

func (writer *Writer) PutInt32BE(v int32) {
	writer.PutUInt32BE(uint32(v))
}

func (writer *Writer) PutInt64LE(v int64) {
	writer.PutUInt64LE(uint64(v))
}

func (writer *Writer) PutInt64BE(v int64) {
	writer.PutUInt64BE(uint64(v))
}
