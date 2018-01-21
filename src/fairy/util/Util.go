package util

import (
	"container/list"
	"os"
	"strings"
)

func SwapList(a *list.List, b *list.List) {
	c := *a
	*a = *b
	*b = c
}

func GetExecName() string {
	var name string
	index := strings.LastIndexAny(os.Args[0], "/\\")
	if index == -1 {
		name = os.Args[0]
	} else {
		name = os.Args[0][index+1:]
	}

	//xxx.exe
	index = strings.LastIndex(name, ".")
	if index == -1 {
		return name
	} else {
		return name[0:index]
	}
}
