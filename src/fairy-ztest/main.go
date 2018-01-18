package main

import (
	// "fairy"
	"fairy"
	"fairy/log"
	"fmt"
)

type Student struct {
	Account string `id:"id"`
}

func Foo() {
	// tmp := make(map[int]int)
	// fmt.Printf("%+v\n", tmp[1])

	// stu := Student{}
	// ptr := &Student{}
	// tstu := reflect.TypeOf(stu)
	// tptr := reflect.TypeOf(ptr)
	// fmt.Printf("ptr:%+v,%+v\n", tptr.Kind(), tptr.Elem().Name())
	// fmt.Printf("stu:%+v,%+v\n", tstu.Kind(), tstu.Name())

	// array := []*Student{}
	// rtype := reflect.TypeOf(array)
	// fmt.Printf("type:%+v,%+v\n", rtype.Kind(), rtype.Elem().Elem().Name())

	// smap := make(map[int]*Student)
	// stype := reflect.TypeOf(smap)
	// fmt.Printf("%+v\n", stype.Key())

	// st := &Student{}
	// rtype := reflect.TypeOf(st)
	// fmt.Printf("ptr:%+v\n", rtype.Name())
	// fmt.Printf("st:%+v\n", reflect.TypeOf(Student{}).Name())
}

func TestTable() {
	data := "STRING\nAccount\nAccount\nJack\naaa"
	records := fairy.ReadTableFromString(data, &Student{})
	for _, cfg := range records.([]*Student) {
		fmt.Printf("%+v\n", cfg)
	}
}

func TestLog() {
	log.GetLogger().Write(log.LEVEL_DEBUG, "%v", 1)
	fairy.Debug("%+v", 1)
}

func main() {
	fmt.Println("start!")
	// test.StartServer()
	// Foo()
	// TestTable()
	TestLog()
	fmt.Println("quit!")
}
