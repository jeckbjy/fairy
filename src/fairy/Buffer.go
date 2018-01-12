package fairy

import (
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

// TODO:impl
type Buffer struct {
	datas    *list.List
	length   int
	position int
	element  *list.Element // 当前游标位置
	offset   int           // 当前位置偏移，element不为空时有效
	mark     int           // 随意位置标识，不做任何校验
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

// 从当前位置分隔成两个
func (self *Buffer) Split(result *Buffer) {
	if self.length == self.position {
		self.Swap(result)
	} else if self.length > 0 {
		// copy
		self.checkCursor()

		// itor := Iterator{}
		// itor.Create(self.datas.Front(), 0, self.position)
		// // 分离
		// for itor.HasNext() {
		// 	data := itor.Next()
		// }

		for iter := self.datas.Front(); iter != nil; iter = iter.Next() {
			// 先移除
			self.datas.Remove(iter)
			//
			data := iter.Value.([]byte)
			if iter == self.element {
				if self.offset >= len(data) {
					result.Append(data)
				} else if self.offset == 0 {
					self.datas.PushFront(data)
				} else {
					// 分成两部分
					front := data[:self.offset]
					back := data[self.offset+1:]
					result.Append(front)
					self.datas.PushFront(back)
				}
				break
			} else {
				result.Append(data)
			}
		}
		//
		self.length -= self.position
		self.position = 0
		self.element = nil
		self.offset = 0
	}
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

// 收缩数据??
// func (self *Buffer) Shrink() {
// }

// 删除当前位置之前的数据
func (self *Buffer) Discard() {
	if self.position == 0 {
		return
	}

	if self.length == self.position {
		self.Clear()
		return
	}

	// copy
	self.checkCursor()
	for iter := self.datas.Front(); iter != nil; iter = iter.Next() {
		if iter == self.element {
			data := iter.Value.([]byte)
			if self.offset >= len(data) {
				self.datas.Remove(iter)
			} else if self.offset > 0 {
				// 删除一部分
				data = data[self.offset:]
				iter.Value = data
			}
			// 停止
			break
		} else {
			self.datas.Remove(iter)
		}
	}

	self.length -= self.position
	self.position = 0
	self.element = nil
	self.offset = 0
}

// 末尾出肉数据,游标移动到末尾
func (self *Buffer) Append(obj interface{}) {
	if obj == nil {
		return
	}

	var data []byte
	switch obj.(type) {
	case []byte:
		data = obj.([]byte)
	case string:
		data = []byte(obj.(string))
	default:
		panic("Buffer.Append bad type!")
	}

	count := len(data)
	if count > 0 {
		self.datas.PushBack(data)
		self.length += count
		self.position = self.length
		self.element = nil
	}
}

// 前边插入数据，游标改为初始位置
func (self *Buffer) Prepend(obj interface{}) {
	if obj == nil {
		return
	}
	var data []byte
	switch obj.(type) {
	case []byte:
		data = obj.([]byte)
	case string:
		data = []byte(obj.(string))
	default:
		panic("Buffer.Prepend bad type!")
	}

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
		pos = offset
	case io.SeekEnd:
		pos = self.length - offset
	}

	if pos < 0 || pos > self.length {
		return fmt.Errorf("seek overflow!")
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
			iter.MoveToEnd()
			// 结束，当前为止
			self.element = iter.element
			self.offset = iter.offset
		} else {
			// 从后向前
		}
	case io.SeekStart:
		iter := Iterator{}
		iter.Create(self.datas.Front(), 0, pos)
		iter.MoveToEnd()
		self.element = iter.element
		self.offset = iter.offset
	case io.SeekEnd:
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

func (self *Buffer) IndexOf(key interface{}) int {
	return self.IndexOfLimit(key, -1)
}

// 从当前位置查找key,最长搜索limit个字节(-1无限制)
func (self *Buffer) IndexOfLimit(key interface{}, limit int) int {
	switch key.(type) {
	case byte:
		// ch := key.(byte)
		for iter := self.datas.Front(); iter != nil; iter = iter.Next() {
			// data := iter.Value.([]byte)
			// find
		}
	case []byte:
	case string:
		//
	default:
	}

	return -1
}

func (self *Buffer) Read(buffer []byte) (int, error) {
	dataNum := len(buffer)
	if self.position+dataNum > self.length {
		return 0, errors.New("Read overflow!")
	}

	self.checkCursor()
	iter := Iterator{}
	iter.Create(self.element, self.offset, dataNum)
	// 拷贝数据
	for iter.HasNext() {
		data := iter.Next()
		copy(buffer[iter.count:], data)
	}

	return dataNum, nil
}

func (self *Buffer) Write(bufffer []byte) (int, error) {
	count := len(bufffer)
	if self.position+count > self.length {
		return 0, errors.New("Write overflow!")
	}

	self.checkCursor()
	iter := Iterator{}
	iter.Create(self.element, self.offset, count)

	for iter.HasNext() {
		data := iter.Next()
		copy(data, bufffer[iter.count:])
	}

	return count, nil
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
	data    []byte
	element *list.Element // 当前节点
	offset  int           // 节点偏移
	length  int           // 需要处理的数据
	count   int           // 已经处理过的数据
}

func (self *Iterator) Create(element *list.Element, offset int, length int) {
	self.element = element
	self.offset = offset
	self.length = length
	if self.length == -1 {
		self.length = math.MaxInt32
	}
	if element != nil {
		self.data = element.Value.([]byte)
	}
}

func (self *Iterator) HasNext() bool {
	return self.element != nil && self.length > 0
}

func (self *Iterator) Next() []byte {
	for self.element != nil {
		if self.offset >= len(self.data) {
			self.element = self.element.Next()
			self.offset = 0
			self.data = nil
			if self.element != nil {
				self.data = self.element.Value.([]byte)
			}
		}
	}

	return self.data
}

func (self *Iterator) NextByte() byte {
	return 0
}

func (self *Iterator) MoveToEnd() {

}
