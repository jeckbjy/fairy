package main

import (
	"fmt"
)

type Student struct {
	Id   int
	Name string
	Desc string
}

func test() {
	// ctype := reflect.TypeOf(Student{})
	// // fmt.Printf("ctype = %+v", ctype)
	// // for i := 0; i < ctype.NumField(); i++ {
	// // 	fmt.Printf("%+v\n", ctype.Field(i))
	// // }

	// // dd := reflect.SliceOf(t)

	// vv := reflect.New(ctype)
	// field := vv.Elem().Field(1)
	// fmt.Printf("aaa:%+v\n", field)
	// field.SetString("Jack")
	// fmt.Printf("bbb:%+v\n", field)
	// fmt.Printf("vv:%+v\n", vv)

	// fmt.Printf("aaaa\n")
	// str := "111|-100,500"
	// tokens := util.SplitNum(str)
	// fmt.Printf("aaa:%+v,%+v\n", tokens, len(tokens))
	// fmt.Printf("aaaa\n")

	// str := "INT,STRING,STRING\nId,Name,Desc\n,,\n1,Jack,Hello"
	// result := []*Student{}
	// records, _ := fairy.ReadTableFromString(str, &result)
	// aa := records.([]*Student)
	// fmt.Printf("%+v", aa)
}

func main() {
	fmt.Println("start!")
	test()
	// test.StartServer()
	fmt.Println("quit!")
}
