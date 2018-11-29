package log

import (
	"fmt"
	"os"
	"path"
)

func NewFileChannel() *FileChannel {
	filename := path.Base(os.Args[0])
	channel := &FileChannel{}
	channel.Init()
	channel.path = fmt.Sprintf("./%+v.log", filename)
	return channel
}

// FileChannel TODO:other strategy:rotate
type FileChannel struct {
	BaseChannel
	path string
	file *os.File
}

// Name return filechannel name
func (f *FileChannel) Name() string {
	return "File"
}

func (f *FileChannel) Write(msg *Message) {
	if f.file == nil {
		return
	}

	output := f.GetOutput(msg)

	f.file.WriteString(output)
}

// Open override base.
func (f *FileChannel) Open() {
	if f.file == nil {
		dir := path.Dir(f.path)
		os.MkdirAll(dir, os.ModePerm)
		var err error
		f.file, err = os.OpenFile(f.path, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			fmt.Printf("log can not open file:%+v\n", f.path)
		}
	}
}

// Close override base.
func (f *FileChannel) Close() {
	if f.file != nil {
		f.file.Close()
	}
}

// SetProperty override base, set config.
func (f *FileChannel) SetProperty(key string, val string) bool {
	if f.BaseChannel.SetProperty(key, val) {
		return true
	}

	switch key {
	case "path":
		f.path = val
		return true
	}

	return false
}
