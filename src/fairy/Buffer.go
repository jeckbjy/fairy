package fairy

import (
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"io"
	"math"
)

func NewBuffer() *Buffer {
	buffer := &Buffer{}
	buffer.datas = list.New()
	return buffer
}

type Buffer struct {
	datas    *list.List    // 数据
	length   int           // 总长度
	position int           // 当前位置
	element  *list.Element // 当前指针
	offset   int           // 当前偏移，element不为空时有效
	mark     int           // 随意位置标识，不做任何校验
}

func (self *Buffer) Front() *list.Element {
	return self.datas.Front()
}

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

func (self *Buffer) Merge(other *Buffer) {
	self.datas.PushBackList(other.datas)
	self.length += other.length
}

// 从当前位置分隔成两个
func (self *Buffer) Split(result *Buffer) {
	if self.length == self.position {
		result.Merge(self)
		self.Clear()
	} else if self.position > 0 {
		// copy
		left := 0
		for iter := self.datas.Front(); iter != nil; {
			temp := iter
			iter = iter.Next()
			//
			data := temp.Value.([]byte)
			size := len(data)
			if size > left {
				result.Append(data[0:left])
				temp.Value = data[left:]
				left = 0
			} else {
				result.Append(data)
				self.datas.Remove(temp)
				left -= size
			}

			if left <= 0 {
				break
			}
		}

		//
		self.length -= self.position
		self.position = 0
		self.element = nil
		self.offset = 0
	}
}

// 删除当前位置之前的数据
func (self *Buffer) Discard() {
	if self.position == 0 {
		return
	}

	if self.length == self.position {
		self.Clear()
		return
	}

	// remove
	left := 0
	for iter := self.datas.Front(); iter != nil; {
		temp := iter
		iter = iter.Next()
		//
		data := temp.Value.([]byte)
		size := len(data)
		if size > left {
			temp.Value = data[left:]
			left = 0
		} else {
			self.datas.Remove(temp)
			left -= size
		}

		if left <= 0 {
			break
		}
	}

	self.length -= self.position
	self.position = 0
	self.element = nil
	self.offset = 0
}

// 合并成一个[]byte
func (self *Buffer) Concat() {
	if self.length == 0 || self.datas.Len() <= 1 {
		return
	}

	data := make([]byte, self.length)
	self.Read(data)
	self.datas.Init()
	self.datas.PushBack(data)
	self.element = self.datas.Front()
	self.offset = self.position
}

// 末尾出肉数据,游标移动到末尾
func (self *Buffer) Append(data []byte) {
	count := len(data)
	if count > 0 {
		self.datas.PushBack(data)
		self.length += count
		self.position = self.length
		self.element = nil
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
func (self *Buffer) Seek(offset int, whence int) error {
	var pos int
	switch whence {
	case io.SeekCurrent:
		pos = self.position + offset
	case io.SeekStart:
		if offset < 0 {
			return fmt.Errorf("Buffer seekstart offset < 0")
		}
		pos = offset
	case io.SeekEnd:
		if offset < 0 {
			return fmt.Errorf("Buffer seekend offset < 0")
		}
		pos = self.length - offset
	}

	if pos < 0 || pos > self.length {
		return fmt.Errorf("buffer seek overflow!")
	}

	if self.position == pos {
		return nil
	}

	if pos == 0 || pos == self.length {
		self.position = pos
		self.element = nil
		return nil
	}

	switch whence {
	case io.SeekCurrent:
		self.checkCursor()
		if offset > 0 {
			// 从前向后
			iter := Iterator{}
			iter.Create(self.element, self.offset, offset)
			iter.MoveEnd()
			// 结束，当前为止
			self.element = iter.element
			self.offset = iter.offset
		} else {
			// 从后向前
			iter := ReverseIterator{}
			iter.Create(self.element, self.offset, -offset)
			iter.MoveEnd()
			self.element = iter.element
			self.offset = iter.offset
		}
	case io.SeekStart:
		iter := Iterator{}
		iter.Create(self.datas.Front(), 0, pos)
		iter.MoveEnd()
		self.element = iter.element
		self.offset = iter.offset
	case io.SeekEnd:
		iter := ReverseIterator{}
		iter.Create(self.element, self.offset, offset)
		iter.MoveEnd()
		self.element = iter.element
		self.offset = iter.offset
	}

	return nil
}

func (self *Buffer) Clear() {
	self.datas.Init()
	self.length = 0
	self.position = 0
	self.element = nil
	self.offset = 0
	self.mark = 0
}

func (self *Buffer) Swap(other *Buffer) {
	temp := Buffer{}
	temp = *other
	*other = *self
	*self = temp
}

// 回到头部
func (self *Buffer) Rewind() {
	self.position = 0
	self.element = self.datas.Front()
	self.offset = 0
}

func (self *Buffer) IndexOf(key string) int {
	return self.IndexOfLimit(key, -1)
}

// 从当前位置查找key,最长搜索limit个字节(-1无限制)
func (self *Buffer) IndexOfLimit(key string, limit int) int {
	self.checkCursor()
	n := len(key)
	switch {
	case n == 0:
		return -1
	case n > self.length-self.position:
		return -1
	case n == 1:
		return self.findByte(key[0], limit)
	}

	// default:
	if limit == -1 {
		limit = math.MaxInt32
	}

	count := 0
	offset := self.offset
	for iter := self.element; iter != nil; iter = iter.Next() {
		data := iter.Value.([]byte)
		for i := offset; i < len(data); i++ {
			if count > limit {
				break
			}
			if self.match(iter, i, key) {
				return self.position + count + i
			}

			count++
		}

		if count > limit {
			break
		}

		offset = 0
	}

	return -1
}

func (self *Buffer) findByte(ch byte, limit int) int {
	iter := Iterator{}
	iter.Create(self.element, self.offset, limit)
	for iter.Next() {
		if pos := bytes.IndexByte(iter.data, ch); pos != -1 {
			return self.position + iter.readNum - len(iter.data) + pos
		}
	}

	return -1
}

func (self *Buffer) match(elem *list.Element, offset int, key string) bool {
	pattern := []byte(key)
	iter := Iterator{}
	iter.Create(self.element, offset, len(pattern))
	for iter.Next() {
		if bytes.Compare(iter.data, pattern) < 0 {
			return false
		}
	}

	return true
}

func (self *Buffer) Read(buffer []byte) (int, error) {
	length := len(buffer)
	if self.position+length > self.length {
		return 0, errors.New("Buffer.Read overflow!")
	}

	self.checkCursor()
	iter := Iterator{}
	iter.Create(self.element, self.offset, length)
	for iter.Next() {
		copy(buffer[iter.readNum:], iter.data)
	}

	// set cursor
	self.element = iter.element
	self.offset = iter.offset
	self.position += length

	return length, nil
}

func (self *Buffer) Write(bufffer []byte) (int, error) {
	self.checkCursor()

	length := len(bufffer)
	count := self.position + length - self.length
	if count > 0 {
		return 0, errors.New("Buffer.write overflow!")
		// resize
		// self.length += count
		// if self.element != nil {
		// 	data := self.element.Value.([]byte)
		// 	remain := cap(data) - len(data)
		// 	var canUse int
		// 	if remain >= count {
		// 		canUse = count
		// 	} else {
		// 		canUse = remain
		// 	}

		//  data = data[:len(data)+canUse]
		// 	self.element.Value = data
		// 	count -= canUse
		// }

		// if count > 0 {
		// 	// create new
		// 	newSize := util.MaxInt(count, 1024)
		// 	data := make([]byte, count, newSize)
		// 	self.datas.PushBack(data)
		// }
	}

	iter := Iterator{}
	iter.Create(self.element, self.offset, length)
	for iter.Next() {
		copy(iter.data, bufffer[iter.readNum:])
	}

	// set cursor
	self.element = iter.element
	self.offset = iter.offset
	self.position += length

	return length, nil
}

func (self *Buffer) SendAll(writer io.Writer) error {
	for iter := self.datas.Front(); iter != nil; iter = iter.Next() {
		data := iter.Value.([]byte)
		if _, err := writer.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func (self *Buffer) ToBytes() []byte {
	if self.length <= 0 {
		return nil
	}
	self.Concat()
	return self.datas.Front().Value.([]byte)
}

// for io.ByteReader
func (self *Buffer) ReadByte() (byte, error) {
	if self.position >= self.length {
		return 0, errors.New("ReadByte overflow!")
	}

	self.checkCursor()

	self.Seek(1, io.SeekCurrent)
	data := self.element.Value.([]byte)
	return data[self.offset], nil
}

func (self *Buffer) checkCursor() {
	if self.element != nil {
		return
	}

	if self.length == 0 {
		return
	}

	if self.position == self.length {
		self.element = self.datas.Back()
		self.offset = len(self.element.Value.([]byte))
	} else {
		self.Seek(self.position, io.SeekStart)
	}
}

////////////////////////////////////////////////////////
//迭代器实现
////////////////////////////////////////////////////////
type Iterator struct {
	element *list.Element // ���前节点
	offset  int           // 当前偏移
	length  int           // 需要读取的总长度
	readNum int           // ���经读取��的
	data    []byte        // 当前读��的数据
}

func (self *Iterator) Create(elem *list.Element, offset int, length int) {
	if length == -1 {
		length = math.MaxInt32
	}

	self.element = elem
	self.offset = offset
	self.length = length
	self.readNum = 0
}

func (self *Iterator) Next() bool {
	if self.element == nil || self.readNum >= self.length {
		return false
	}

	needRead := self.length - self.readNum
	for {
		data := self.element.Value.([]byte)
		left := len(data) - self.offset
		if left <= 0 {
			// ignore zero data
			self.element = self.element.Next()
			self.offset = 0
			if self.element == nil {
				return false
			}
			continue
		}

		// read data
		if left > needRead {
			self.data = data[self.offset : self.offset+needRead]
			self.offset += needRead
			self.readNum += needRead
		} else {
			self.data = data[self.offset:]
			self.element = self.element.Next()
			self.offset = 0
			self.readNum += left
		}
		break
	}

	return true
}

func (self *Iterator) MoveEnd() {
	needRead := self.readNum
	for needRead > 0 && self.element != nil {
		data := self.element.Value.([]byte)
		left := len(data) - self.offset
		if left > needRead {
			self.offset += needRead
			needRead = 0
		} else {
			needRead -= left
			self.element = self.element.Next()
			self.offset = 0
		}
	}
}

////////////////////////////////////////////////////////
//反向迭代器
////////////////////////////////////////////////////////
type ReverseIterator struct {
	element *list.Element // 当前节点
	offset  int           // 当前偏移
	length  int           // 需要读取的总长度
	readNum int           // 已经读取完的
	data    []byte        // 当前读取的数据
}

func (self *ReverseIterator) Create(elem *list.Element, offset int, length int) {
	self.element = elem
	self.offset = offset
	self.length = length
	self.readNum = 0
}

func (self *ReverseIterator) Next() bool {
	if self.element == nil || self.readNum >= self.length {
		return false
	}

	needRead := self.length - self.readNum
	for {
		if self.offset == 0 {
			self.element = self.element.Next()
			self.offset = 0
			if self.element == nil {
				return false
			}
			continue
		}

		data := self.element.Value.([]byte)
		if self.offset == -1 {
			self.offset = len(data)
		}

		if self.offset > needRead {
			start := self.offset - needRead
			self.data = data[start:self.offset]
			self.readNum += needRead
		} else {
			self.data = data[0:self.offset]
			self.element = self.element.Prev()
			self.offset = -1
			self.readNum += self.offset
		}

		break
	}

	return true
}

func (self *ReverseIterator) MoveEnd() {
	needRead := self.readNum
	for needRead > 0 && self.element != nil {
		data := self.element.Value.([]byte)
		if self.offset == -1 {
			self.offset = len(data)
		}
		if self.offset >= needRead {
			self.offset -= needRead
			needRead = 0
		} else {
			self.element = self.element.Prev()
			self.offset = -1
			needRead -= self.offset
		}
	}
}
