package main

import (
	"fairy-ztest/echo"
	"flag"
	"fmt"
)

func main() {
	fmt.Println("start!")
	mode := flag.String("m", "client", "server mode")
	flag.Parse()
	if *mode == "server" {
		echo.StartServer()
	} else {
		echo.StartClient()
	}
	fmt.Println("quit!")
}
