package main

import (
	"fairy-ztest/test"
	"fmt"
)

func main() {
	fmt.Println("start!")
	test.StartServer()
	fmt.Println("quit!")
}
