package netlisten_test

import "testing"
import "chatroom/netlisten"
import "fmt"

func TestNewNetListener(t *testing.T) {
	listener := netlisten.NewNetListener("tcp://0.0.0.0", "9999")
	fmt.Println(listener)
}

func TestRun(t *testing.T) {
	listener := netlisten.NewNetListener("tcp://0.0.0.0", "9999")
	listener.Run()
}
