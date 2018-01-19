package log

import (
	"os"
)

func NewFileChannel() *FileChannel {
	channel := &FileChannel{}
	return channel
}

type FileChannel struct {
	BaseChannel
	path string
	file *os.File
}

func (self *FileChannel) Name() string {
	return "File"
}

func (self *FileChannel) Write(msg *Message) {
}

func (self *FileChannel) Open() {
	if self.file == nil {
		self.file, _ = os.Open(self.path)
	}
}

func (self *FileChannel) Close() {

}
