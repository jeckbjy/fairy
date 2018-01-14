package log

import (
	"os"
)

type FileChannel struct {
	BaseChannel
	path string
	file *os.File
}

func (self *FileChannel) Write(msg *Message) {
	//
	if self.CanWrite(msg.Level) {
		//
	}
}

func (self *FileChannel) Open() {
	os.Open(self.path)
}

func (self *FileChannel) Close() {

}