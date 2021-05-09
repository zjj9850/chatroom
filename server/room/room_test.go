package room_test

import "testing"
import "chatroom/room"
import "fmt"

func TestNewRoomMgr(t *testing.T) {
	room.NewRoomMgr(nil, nil)
}

func TestInit(t *testing.T) {
	mgr := room.NewRoomMgr(nil, nil)
	mgr.Init()
	fmt.Println(mgr)
}

func TestRun(t *testing.T) {
	mgr := room.NewRoomMgr(nil, nil)
	mgr.Init()
	mgr.Run()
}
