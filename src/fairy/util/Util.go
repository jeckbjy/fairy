package util

import "container/list"

func SwapList(a *list.List, b *list.List) {
	c := *a
	*a = *b
	*b = c
}
