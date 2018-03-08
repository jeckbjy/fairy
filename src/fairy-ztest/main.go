package main

import (
	"fmt"
)

func test_slice() {
	aa := make([]byte, 0, 5)
	aa = append(aa, '1')
	aa = append(aa, '2')
	fmt.Printf("len=%+v\n", len(aa))
	bb := aa[2:]
	fmt.Printf("len.a=%+v,len.b=%+v\n", len(aa), len(bb))
}

func main() {
	test_slice()
	// echo.Test()
	// chat.Test()
}
