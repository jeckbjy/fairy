package util

import (
	"strconv"
)

type PatternAction struct {
	Key      rune
	Prepend  string // xxx%
	Property string // %[name]
}

type PatternActionArray []*PatternAction

func (self *PatternAction) PropInt() int {
	val, _ := strconv.Atoi(self.Property)
	return val
}

func (self *PatternAction) PropStr() string {
	return self.Property
}

// 通用解析规则%?[prop],example:[%y-%m-%d %H:%M:%S][%q][%U:%u][%t]
func ParsePattern(format string) []*PatternAction {
	actions := []*PatternAction{}

	end := len(format)
	cur := 0
	for cur < end {
		act := &PatternAction{}
		// parse prepend
		for beg := cur; ; cur++ {
			if cur >= end || format[cur] == '%' {
				if beg < cur {
					act.Prepend = format[beg:cur]
				}
				break
			}
		}

		// check end
		if cur == end {
			actions = append(actions, act)
			break
		}
		// parse key
		cur++
		act.Key = rune(format[cur])
		cur++
		// parse property
		if cur < end && format[cur] == '[' {
			cur++
			for beg := cur; cur < end ; cur++ {
				if format[cur] == ']' {
					act.Property = format[beg:cur]
					cur++
					break
				}
			}
		}
		actions = append(actions, act)
	}

	return actions
}