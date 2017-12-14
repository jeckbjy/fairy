package fairy

import (
	"container/list"
	"errors"
	"io"
)

func NewBuffer() *Buffer {
	buffer := &Buffer{}
	buffer.datas = list.New()
	return buffer
}

//迭代器实现
type BufferIterator struct {
	element *list.Element
	offset  int
}

// todo:Iterator
type Buffer struct {
	datas    *list.List
	length   int
	position int
	mark     int           // 随意位置标识，不做任何校验
	element  *list.Element // 当前游标位置
	offset   int
}

// 从当前位置分隔成两个
func (self *Buffer) Split(buffer *Buffer) error {
	return nil
}

// 合并成一个[]byte
func (self *Buffer) Concat() {

}

// 收缩数据??
// func (self *Buffer) Shrink() {
// }

// 删除当前位置之前的数据
func (self *Buffer) Discard() {

}

func (self *Buffer) DiscardCount(count int) {
	self.Seek(count, io.SeekCurrent)
	self.Discard()
}

// 末尾出肉数据,不修改游标
func (self *Buffer) Append(data []byte) {
	leng := len(data)
	if leng > 0 {
		self.datas.PushBack(data)
		self.length += leng
	}
}

// 前边插入数据，游标改为初始位置
func (self *Buffer) Prepend(data []byte) {
	count := len(data)
	if count > 0 {
		self.datas.PushFront(data)
		self.length += count
		self.element = self.datas.Front()
		self.offset = 0
		self.position = 0
	}
}

// io.SeekStart
func (self *Buffer) Seek(offset int, whence int) {
	switch whence {
	case io.SeekCurrent:
	case io.SeekStart:
	case io.SeekEnd:
	}
}

// 回到头部
func (self *Buffer) Rewind() {
	self.position = 0
	self.element = self.datas.Front()
	self.offset = 0
}

func (self *Buffer) IndexOf(key interface{}) int {
	return self.IndexOfLimit(key, -1)
}

// 从当前位置查找key,最长搜索limit个字节(-1无限制)
func (self *Buffer) IndexOfLimit(key interface{}, limit int) int {
	// 比较
	if _, ok := key.(byte); ok {
		// 从当前位置查找一个字符
	} else if _, ok := key.([]byte); ok {
		// 从当前位置查找多个字符
	} else if _, ok := key.(string); ok {
		// 字符
	}
	return -1
}

// func (self *Buffer) LastIndexOf(key interface{}) int {
// 	return -1
// }

func (self *Buffer) HasRemain(count int) bool {
	return self.length-self.position >= count
}

func (self *Buffer) Empty() bool {
	return self.length == 0
}

func (self *Buffer) Length() int {
	return self.length
}

func (self *Buffer) Position() int {
	return self.position
}

func (self *Buffer) Eof() bool {
	return self.position >= self.length
}

func (self *Buffer) SetMark(value int) {
	self.mark = value
}

func (self *Buffer) GetMark() int {
	return self.mark
}

func (self *Buffer) Read(buffer []byte) (int, error) {
	dataNum := len(buffer)
	if self.position+dataNum > self.length {
		return 0, errors.New("Read overflow!")
	}

	self.travel(dataNum, func(offset int, value []byte) {
		copy(buffer[offset:], value)
	})

	return dataNum, nil
}

func (self *Buffer) Write(bufffer []byte) (int, error) {
	dataNum := len(bufffer)
	if self.position+dataNum > self.length {
		return 0, errors.New("Write overflow!")
	}

	self.travel(dataNum, func(offset int, value []byte) {
		copy(value, bufffer[offset:])
	})

	return dataNum, nil
}

func (self *Buffer) travel(count int, cb func(int, []byte)) int {
	needNum := count
	dataOff := 0
	copyNum := 0
	elem := self.element
	offs := self.offset
	for elem != nil {
		value := elem.Value.([]byte)
		leftNum := len(value) - offs
		if leftNum > 0 {
			if leftNum >= needNum {
				copyNum = needNum
			} else {
				copyNum = leftNum
			}
			cb(dataOff, value[offs:offs+copyNum])
			dataOff += copyNum
			needNum -= copyNum
			offs += copyNum
			if needNum <= 0 {
				break
			}
		}

		elem = elem.Next()
		offs = 0
	}

	// 修改当前游标，为nil如何处理？
	if elem == nil {
		// 设置为最后一个?
		elem = self.datas.Back()
		offs = len(elem.Value.([]byte))
	}

	self.element = elem
	self.offset = offs
	self.position += dataOff

	// 返回处理了多少字节
	return dataOff
}

func (self *Buffer) ToBytes() []byte {
	return nil
}

func (self *Buffer) ReadByte() (byte, error) {
	if self.position >= self.length {
		return 0, errors.New("ReadByte overflow!")
	}

	self.Seek(1, io.SeekCurrent)
	return self.getCurrentChunk()[self.offset], nil
}

// 判断当前游标位置是否是结尾
func (self *Buffer) checkEnd() []byte {
	buffer := self.element.Value.([]byte)
	if self.offset < len(buffer) {
		return buffer
	}

	self.element = self.element.Next()
	self.offset = 0
	buffer = self.element.Value.([]byte)
	return buffer
}

func (self *Buffer) getCurrentChunk() []byte {
	buffer := self.element.Value.([]byte)
	return buffer
}
