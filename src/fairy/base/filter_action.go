package base

import "fairy"

const (
	FILTER_ACTION_STOP  = 0 // 执行终止,比如消息包没有解析完
	FILTER_ACTION_NEXT  = 1 // 执行下一个，大部分情况
	FILTER_ACTION_LAST  = 2 // 执行最后一个
	FILTER_ACTION_FIRST = 3 // 执行第一个
)

var (
	gStopAction  = NewFilterAction(FILTER_ACTION_STOP)
	gNextAction  = NewFilterAction(FILTER_ACTION_NEXT)
	gLastAction  = NewFilterAction(FILTER_ACTION_LAST)
	gFirstAction = NewFilterAction(FILTER_ACTION_FIRST)
)

func NewFilterAction(ftype int) fairy.FilterAction {
	action := &FilterAction{}
	action.ftype = ftype
	return action
}

type FilterAction struct {
	ftype int
}

func (self *FilterAction) Type() int {
	return self.ftype
}
