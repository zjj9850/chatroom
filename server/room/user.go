package room

import (
	"github.com/panjf2000/gnet"
)

const GM_PASSWORD = "chatroom"

type User struct {
	Name      string
	IsGM      bool
	LoginTime int64
	RoomId    uint32
	Conn      gnet.Conn
}

func newUser(name string, pwd string) *User {
	return &User{
		Name:   name,
		IsGM:   pwd == GM_PASSWORD,
		RoomId: 0,
	}
}
