package main

import (
	"fairy-ztest/echo"
	"flag"
)

func main() {
	// fmt.Println("start!")
	mode := flag.String("m", "server", "server mode")
	flag.Parse()
	if *mode == "server" {
		echo.StartServer()
	} else {
		echo.StartClient()
	}
	// fmt.Println("quit!")
}
