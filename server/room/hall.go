package room

import (
	"chatroom/logkit"
	"chatroom/netlisten"
	"github.com/panjf2000/gnet"
	"sync"
	"time"
)

type Hall struct {
	ConnMap          sync.Map
	UserMapByName    sync.Map
	UserNameByConnId sync.Map
}

func newHall() *Hall {
	return &Hall{}
}

func (self *Hall) getUserByName(name string) *User {
	if iUser, _ := self.UserMapByName.Load(name); iUser != nil {
		user := iUser.(*User)
		return user
	}
	return nil
}

func (self *Hall) getUserByConnId(connId uint32) *User {
	if n, _ := self.UserNameByConnId.Load(connId); n != nil {
		userName := n.(string)
		if iUser, _ := self.UserMapByName.Load(userName); iUser != nil {
			user := iUser.(*User)
			return user
		}
	}
	return nil
}

func (self *Hall) login(userName string, pwd string, connId uint32) {
	conn, _ := self.ConnMap.Load(connId)
	if conn != nil {
		user := newUser(userName, pwd)
		user.LoginTime = time.Now().Unix()
		user.RoomId = 0
		user.Conn = conn.(gnet.Conn)
		self.UserMapByName.Store(userName, user)
		self.UserNameByConnId.Store(connId, userName)
	}
}

func (self *Hall) run(server *netlisten.NetListener, leaveChan chan<- *User) {
	for {
		select {
		case conn := <-server.OpenChannel:
			connId := conn.Context().(uint32)
			self.ConnMap.Store(connId, conn)
			logkit.Infof("%s ,ConnId:%d Connect the chathall", conn.RemoteAddr().String(), connId)
		case connId := <-server.CloseChannel:
			self.ConnMap.Delete(connId)
			if n, _ := self.UserNameByConnId.Load(connId); n != nil {
				userName := n.(string)
				self.UserNameByConnId.Delete(connId)
				if iUser, _ := self.UserMapByName.Load(userName); iUser != nil {
					user := iUser.(*User)
					leaveChan <- user
					self.UserMapByName.Delete(userName)
				}
			}
			logkit.Infof("ConnId:%d Disconnect the chathall", connId)
		}
	}
}
