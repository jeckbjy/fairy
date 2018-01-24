package echo

import (
	"fairy"
	"testing"
)

func EchoTest(t *testing.T) {
	go StartServer()
	go StartClient()
	fairy.WaitExit()
}
