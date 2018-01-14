package log

import (
	"container/list"
)

type MessageQueue struct {
	messages *list.List
}