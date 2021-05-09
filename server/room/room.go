package room

import (
	"chatroom/protocol"
	"sync"
)

const MAX_RECENT_CHAT_MSG = 50

type Room struct {
	RoomId        uint32
	UserMapByName sync.Map
	RecentMsg     []*protocol.ChatRes
}

func newRoom(roomId uint32) *Room {
	return &Room{
		RoomId:    roomId,
		RecentMsg: make([]*protocol.ChatRes, 0, MAX_RECENT_CHAT_MSG),
	}
}

func (self *Room) pushMsg(chatRes *protocol.ChatRes) {
	if len(self.RecentMsg) >= MAX_RECENT_CHAT_MSG {
		self.RecentMsg = self.RecentMsg[1:]
	}
	self.RecentMsg = append(self.RecentMsg, chatRes)
}

func (self *Room) getUserByName(name string) *User {
	if iUser, _ := self.UserMapByName.Load(name); iUser != nil {
		user := iUser.(*User)
		return user
	}
	return nil
}

func (self *Room) joinUser(user *User) {
	self.UserMapByName.Store(user.Name, user)
}

func (self *Room) removeUserByName(name string) {
	self.UserMapByName.Delete(name)
}
