package test

import (
	"fairy"
	"io"
	"testing"
)

func NewBuffer() *fairy.Buffer {
	buffer := fairy.NewBuffer()
	buffer.Append([]byte("key"))
	buffer.Append([]byte(":"))
	buffer.Append([]byte("val"))
	buffer.Seek(0, io.SeekStart)
	return buffer
}

func TestBuffer(t *testing.T) {
	buffer := NewBuffer()

	pos := buffer.IndexOf(":")
	if pos != 3 {
		t.Errorf("buffer indexof fail:%+v", pos)
	}

	buffer.Seek(0, io.SeekEnd)
	if buffer.Position() != buffer.Length() {
		t.Errorf("buffer seek end fail:%+v, %+v", buffer.Position(), buffer.Length())
	}

	buffer.Seek(-1, io.SeekEnd)
	if buffer.Position() != buffer.Length()-1 {
		t.Errorf("buffer seek fail!:%+v,%+v", buffer.Position(), buffer.Length())
	}

	buffer.Seek(4, io.SeekStart)
	if buffer.Position() != 4 {
		t.Error("buffer seek fail!")
	}

	buffer.Seek(-1, io.SeekCurrent)
	if buffer.Position() != 3 {
		t.Error("buffer seek fail")
	}

	data := make([]byte, 3)
	buffer.Seek(4, io.SeekStart)
	buffer.Read(data)
	if string(data) != "val" {
		t.Error("buffer read fail")
	}

	buffer.Seek(2, io.SeekStart)
	buffer.Write([]byte("a?b"))
	if buffer.String() != "kea?bal" {
		t.Error("buffer write fail")
	}

	buffer = NewBuffer()
	buffer.Seek(3, io.SeekStart)
	buffer.Discard()
	if buffer.String() != ":val" {
		t.Error("buffer discard fail!")
	}

	buffer = NewBuffer()
	b1 := fairy.NewBuffer()
	buffer.Seek(3, io.SeekStart)
	buffer.Split(b1)
	if b1.String() != "key" || buffer.String() != ":val" {
		t.Error("split fail")
	}
}
