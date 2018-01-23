package log

import (
	"fairy/util"
	"fmt"
	"os"
	"path"
)

func NewFileChannel() *FileChannel {
	channel := &FileChannel{}
	channel.Init()
	channel.path = fmt.Sprintf("./%+v.log", util.GetExecName())
	return channel
}

// TODO:other strategy:rotate
type FileChannel struct {
	BaseChannel
	path string
	file *os.File
}

func (self *FileChannel) Name() string {
	return "File"
}

func (self *FileChannel) Write(msg *Message) {
	output := self.GetOutput(msg)
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Printf("%+v,%+v\n", output, err)
	// 	}
	// }()

	self.file.WriteString(output)
}

func (self *FileChannel) Open() {
	if self.file == nil {
		dir := path.Dir(self.path)
		os.MkdirAll(dir, os.ModePerm)
		var err error
		self.file, err = os.OpenFile(self.path, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Printf("log can not open file:%+v\n", self.path)
		}
	}
}

func (self *FileChannel) Close() {
	if self.file != nil {
		self.file.Close()
	}
}
