package chat

import (
	"fairy/util"
	"testing"
)

func TestChat(t *testing.T) {
	StartServer()
	StartClient()
	util.Sleep(20 * 1000)
}
