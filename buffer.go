package fairy

import (
	"bytes"
	"container/list"
	"fmt"
	"io"
	"math"
)

const defaultGrowSize = 1024

var errNotFindLine = fmt.Errorf("err not find line!")

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

// GetSpace = back of (capacity - length)
func (self *Buffer) GetSpace() []byte {
	if self.datas.Len() == 0 {
		return nil
	}

	data := self.datas.Back().Value.([]byte)
	caps := cap(data)
	leng := len(data)
	if leng == caps {
		return nil
	}

	return data[leng:caps]
}

// ExtendSpace 将最后一个扩展count个字节,配合GetSpace使用
func (self *Buffer) ExtendSpace(count int) error {
	if self.datas.Len() == 0 {
		return fmt.Errorf("buffer extend space, but back data")
	}

	back := self.datas.Back()
	data := back.Value.([]byte)
	leng := len(data)
	back.Value = data[0 : leng+count]

	self.length += count
	self.position = self.length
	self.element = nil
	return nil
}

// Merge 合并两个buffer为一个
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
		left := self.position
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

	result.Seek(0, io.SeekStart)
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
	left := self.position
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

	data := make([]byte, 0, self.length)
	for iter := self.datas.Front(); iter != nil; iter = iter.Next() {
		data = append(data, iter.Value.([]byte)...)
	}
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

// Seek 移动当前位置
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
		if offset > 0 {
			offset = -offset
		}
		pos = self.length + offset
	}

	if pos < 0 || pos > self.length {
		return fmt.Errorf("buffer seek overflow")
	}

	// check element
	if self.element == nil {
		if pos == 0 {
			self.element = self.datas.Front()
			self.offset = 0
			self.position = 0
			return nil
		}

		if pos == self.length {
			self.element = self.datas.Back()
			self.offset = len(self.element.Value.([]byte))
			self.position = pos
			return nil
		}

		if whence == io.SeekCurrent {
			whence = io.SeekStart
		}
	} else if self.position == pos {
		// not need
		return nil
	}

	switch whence {
	case io.SeekCurrent:
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

	self.position = pos

	return nil
}

// Clear 清空数据
func (self *Buffer) Clear() {
	self.datas.Init()
	self.length = 0
	self.position = 0
	self.element = nil
	self.offset = 0
	self.mark = 0
}

// Swap 交换两个buffer
func (self *Buffer) Swap(other *Buffer) {
	temp := Buffer{}
	temp = *other
	*other = *self
	*self = temp
}

// Rewind 回到头部
func (self *Buffer) Rewind() {
	self.position = 0
	self.element = self.datas.Front()
	self.offset = 0
}

// IndexOf 从当前位置开始查询字符串位置
func (self *Buffer) IndexOf(key string) int {
	return self.IndexOfLimit(key, -1)
}

// IndexOfLimit 从当前位置查找key,最长搜索limit个字节(-1无限制)
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

// Peek 获取数据但不修改当前位置
func (self *Buffer) Peek(buffer []byte) (int, error) {
	length := len(buffer)
	if self.position+length > self.length {
		return 0, fmt.Errorf("Buffer.Read overflow!")
	}

	self.checkCursor()
	iter := Iterator{}
	iter.Create(self.element, self.offset, length)
	for iter.Next() {
		copy(buffer[iter.lastNum:], iter.data)
	}

	return length, nil
}

// Read 读取数据并移动当前位置
func (self *Buffer) Read(buffer []byte) (int, error) {
	length := len(buffer)
	if self.position+length > self.length {
		return 0, fmt.Errorf("Buffer.Read overflow!")
	}

	self.checkCursor()
	iter := Iterator{}
	iter.Create(self.element, self.offset, length)
	for iter.Next() {
		copy(buffer[iter.lastNum:], iter.data)
	}

	// set cursor
	self.element = iter.element
	self.offset = iter.offset
	self.position += length

	return length, nil
}

func (self *Buffer) grow(count int) {
	// check empty
	if self.datas.Len() == 0 {
		newCap := count + count>>1
		if newCap < defaultGrowSize {
			newCap = defaultGrowSize
		}
		data := make([]byte, count, newCap)
		self.datas.PushBack(data)
	} else {
		// expand last first
		back := self.datas.Back()
		data := back.Value.([]byte)
		caps := cap(data)
		leng := len(data)
		left := caps - leng

		// back enough
		if left >= count {
			back.Value = data[0 : leng+count]
		} else {
			// step1: resize left
			if left > 0 {
				back.Value = data[0:caps]
			}

			// step2: alloc new
			newCount := count - left
			newCap := newCount + newCount>>1
			if newCap < defaultGrowSize {
				newCap = defaultGrowSize
			}
			data = make([]byte, newCount, newCap)
			self.datas.PushBack(data)
		}
	}

	// resize
	self.length += count
}

// Write 实现io.Writer接口
func (self *Buffer) Write(bufffer []byte) (int, error) {
	self.checkCursor()

	length := len(bufffer)
	count := self.position + length - self.length
	if count > 0 {
		self.grow(count)
		// return 0, fmt.Errorf("Buffer.write overflow!")
	}

	iter := Iterator{}
	iter.Create(self.element, self.offset, length)
	for iter.Next() {
		copy(iter.data, bufffer[iter.lastNum:])
	}

	// set cursor
	self.element = iter.element
	self.offset = iter.offset
	self.position += length

	return length, nil
}

func (buffer *Buffer) ReadAll(reader io.Reader) error {
	// 先尝试填充末尾
	if buffer.datas.Len() > 0 {
		back := buffer.datas.Back()
		data := back.Value.([]byte)
		caps := cap(data)
		leng := len(data)
		if leng < caps {
			p := data[leng:caps]
			count, err := reader.Read(p)
			if count > 0 {
				back.Value = data[0 : leng+count]
				buffer.length += count
			}

			if err != nil {
				return err
			}
		}
	}

	// 继续读取
	for {
		p := make([]byte, defaultGrowSize)
		count, err := reader.Read(p)
		if count > 0 {
			buffer.Append(p[:count])
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (buffer *Buffer) WriteAll(writer io.Writer) error {
	for iter := buffer.Front(); iter != nil; iter = iter.Next() {
		data := iter.Value.([]byte)
		_, err := writer.Write(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Buffer) Bytes() []byte {
	if self.length <= 0 {
		return nil
	}
	self.Concat()
	return self.datas.Front().Value.([]byte)
}

func (self *Buffer) String() string {
	if self.length <= 0 {
		return ""
	}
	self.Concat()
	data := self.datas.Front().Value.([]byte)
	return string(data)
}

func (self *Buffer) ReadToEnd() []byte {
	if self.length == 0 {
		return nil
	}

	self.Concat()
	data := self.datas.Front().Value.([]byte)
	if self.position != 0 {
		return data[self.position:]
	}

	return data
}

// ReadByte 实现接口io.ByteReader
func (self *Buffer) ReadByte() (byte, error) {
	if self.position >= self.length {
		return 0, fmt.Errorf("ReadByte overflow!")
	}

	self.checkCursor()

	result := self.element.Value.([]byte)[self.offset]
	self.Seek(1, io.SeekCurrent)
	return result, nil
}

// ReadUntil 读取到key位置
func (self *Buffer) ReadUntil(key byte) (string, error) {
	if self.position == self.length {
		return "", fmt.Errorf("buffer end,cannot find key=%+v", string(key))
	}

	self.checkCursor()
	index := self.findByte(key, -1)
	if index == -1 {
		return "", fmt.Errorf("cannot find key=%+v", key)
	}

	data := make([]byte, index)
	self.Read(data)
	// remove key
	self.Seek(1, io.SeekCurrent)
	return string(data), nil
}

// ReadLine 读取到\n或\r\n为止
func (self *Buffer) ReadLine() (*Buffer, error) {
	delimiter := 1

	pos := self.IndexOf("\n")
	if pos == -1 {
		return nil, errNotFindLine
	}

	if pos > 0 {
		self.Seek(-1, io.SeekCurrent)
		if ch, _ := self.ReadByte(); ch == '\r' {
			delimiter = 2
			pos--
			self.Seek(-1, io.SeekCurrent)
		}
	}

	result := NewBuffer()
	self.Seek(pos, io.SeekStart)
	self.Split(result)

	self.Seek(delimiter, io.SeekStart)
	self.Discard()

	return result, nil
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

// 遍历整个数组
func (self *Buffer) Visit(cb func([]byte) bool) {
	for iter := self.datas.Front(); iter != nil; iter = iter.Next() {
		data := iter.Value.([]byte)
		if !cb(data) {
			break
		}
	}
}

////////////////////////////////////////////////////////
//迭代器实现
////////////////////////////////////////////////////////
type Iterator struct {
	element *list.Element // 当前节点
	offset  int           // 当前偏移
	length  int           // 需要读取的总长度
	lastNum int           // 上次读取位置
	readNum int           // 已经读取的长度
	data    []byte        // 当前读取的数据
}

func (self *Iterator) Create(elem *list.Element, offset int, length int) {
	if length == -1 {
		length = math.MaxInt32
	}

	self.element = elem
	self.offset = offset
	self.length = length
	self.readNum = 0
	self.lastNum = 0
}

func (self *Iterator) Next() bool {
	if self.element == nil || self.readNum >= self.length {
		return false
	}

	self.lastNum = self.readNum

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
	needRead := self.length - self.readNum
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
	lastNum int           // 上次读取位置
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

	self.lastNum = self.readNum

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
