package util

import (
	"container/list"
	"os"
	"strings"
	"sync"
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

func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// 非递归锁不能嵌套使用
var gOnceMutex sync.Mutex

func Once(inst interface{}, cb func()) {
	if inst == nil {
		gOnceMutex.Lock()
		if inst == nil {
			cb()
		}
		gOnceMutex.Unlock()
	}
}

// 防止递归锁
var gOnceMutexEx sync.Mutex

func OnceEx(inst interface{}, cb func()) {
	if inst == nil {
		gOnceMutexEx.Lock()
		if inst == nil {
			cb()
		}
		gOnceMutexEx.Unlock()
	}
}
