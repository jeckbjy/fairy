package main

import (
	"fairy-ztest/test"
	"fmt"
	"reflect"
)

type Student struct {
	Account string
}

func Foo() {
	st := &Student{}
	rtype := reflect.TypeOf(st)
	fmt.Printf("ptr:%+v\n", rtype.Name())
	fmt.Printf("st:%+v\n", reflect.TypeOf(Student{}).Name())
}

func main() {
	fmt.Println("start!")
	test.StartServer()
	// Foo()
	fmt.Println("quit!")
}
