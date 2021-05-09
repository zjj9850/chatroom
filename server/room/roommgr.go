package room

import (
	"chatroom/filter"
	"chatroom/logkit"
	"chatroom/netlisten"
	"chatroom/protocol"
	"fmt"
	"github.com/panjf2000/gnet"
	"google.golang.org/protobuf/proto"
	"strconv"
	"strings"
	"sync"
	"time"
)

type MsgHandlerFunc = func([]byte, uint32)

type RoomMgr struct {
	server         *netlisten.NetListener
	chatHall       *Hall
	roomMap        sync.Map
	wordFilter     *filter.WordFilter
	mapHandler     map[string]MsgHandlerFunc
	UserLeaveChan  chan *User
	WordStatistics map[int64]map[string]int
	wordStatChan   chan string
}

func NewRoomMgr(listener *netlisten.NetListener, wordFilter *filter.WordFilter) *RoomMgr {
	return &RoomMgr{
		server:         listener,
		chatHall:       newHall(),
		wordFilter:     wordFilter,
		mapHandler:     make(map[string]MsgHandlerFunc),
		UserLeaveChan:  make(chan *User),
		WordStatistics: make(map[int64]map[string]int),
		wordStatChan:   make(chan string),
	}
}

func (self *RoomMgr) registerHandler(msg proto.Message, f MsgHandlerFunc) {
	self.mapHandler[string(msg.ProtoReflect().Descriptor().FullName())] = f
}

func (self *RoomMgr) Init() {
	self.registerHandler(&protocol.LoginReq{}, self.HandleLogin)
	self.registerHandler(&protocol.JoinRoomReq{}, self.HandleJoinRoom)
	self.registerHandler(&protocol.ChatReq{}, self.HandleChat)
	self.registerHandler(&protocol.PrivateChatReq{}, self.HandlePrivateChat)
}

func (self *RoomMgr) Run() {
	go self.chatHall.Run(self.server, self.UserLeaveChan)
	go self.NetFrameCheck()
	go self.UserLeaveCheck()
	go self.WordStat()
}

func (self *RoomMgr) RemoveUser(roomId uint32, userName string) {
	if iRoom, _ := self.roomMap.Load(roomId); iRoom != nil {
		room := iRoom.(*Room)
		room.RemoveUserByName(userName)
		self.NotifyRoom(room, fmt.Sprintf("User %s has leaved No.%d Room", userName, roomId), userName)
	}
}

func (self *RoomMgr) WordStat() {
	for {
		select {
		case content := <-self.wordStatChan:
			nowSec := time.Now().Unix()
			statMap, e := self.WordStatistics[nowSec]
			if !e {
				statMap = make(map[string]int)
			}
			wordList := strings.Split(content, " ")
			for _, word := range wordList {
				if c, we := statMap[word]; we {
					statMap[word] = c + 1
				} else {
					statMap[word] = 1
				}
			}
			self.WordStatistics[nowSec] = statMap

			deleteSec := make([]int64, 0)
			for sec, _ := range self.WordStatistics {
				if nowSec-sec > 60 {
					deleteSec = append(deleteSec, sec)
				}
			}

			for _, sec := range deleteSec {
				delete(self.WordStatistics, sec)
			}
		}
	}
}

func (self *RoomMgr) GetTopWord(sec int) (string, int) {
	if sec > 60 || sec < 1 {
		return "", 0
	}

	nowSec := time.Now().Unix()
	wordMap := make(map[string]int)

	for s := nowSec - int64(sec); s <= nowSec; s++ {
		statMap, e := self.WordStatistics[s]
		if e {
			for str, count := range statMap {
				if _, we := wordMap[str]; we {
					wordMap[str] += count
				} else {
					wordMap[str] = count
				}
			}
		}
	}

	maxCount := 0
	maxStr := ""
	for str, count := range wordMap {
		if count >= maxCount {
			maxStr = str
			maxCount = count
		}
	}

	return maxStr, maxCount
}

func (self *RoomMgr) UserLeaveCheck() {
	for {
		select {
		case user := <-self.UserLeaveChan:
			self.RemoveUser(user.RoomId, user.Name)
		}
	}
}

func (self *RoomMgr) NetFrameCheck() {
	for frameInfo := range self.server.FrameChannel {
		connId := frameInfo.ConnId

		msg := &protocol.Message{}
		err := proto.Unmarshal(frameInfo.Frame, msg)
		if err != nil {
			logkit.Errorf("Unmarshal Message failed,Err:%s,ConnId:%d", err.Error(), connId)
			continue
		}

		handler, ok := self.mapHandler[msg.Type]
		if ok {
			handler(msg.Data, connId)
		} else {
			logkit.Warning("UnRegister msg recv:", msg.Type)
		}
	}
}

func (self *RoomMgr) SendMsg(conn gnet.Conn, msg proto.Message) {
	if conn == nil {
		return
	}

	packed := &protocol.Message{}
	packed.Type = string(msg.ProtoReflect().Descriptor().FullName())
	packed.Data, _ = proto.Marshal(msg)
	send_packed, err := proto.Marshal(packed)
	if err != nil {
		logkit.Error("Marshal Failed", err.Error())
	} else {
		self.server.SendMsg(conn, send_packed)
	}
}

func (self *RoomMgr) HandleLogin(data []byte, connId uint32) {
	iConn, _ := self.chatHall.ConnMap.Load(connId)
	if iConn == nil {
		return
	}

	conn := iConn.(gnet.Conn)

	msg := &protocol.LoginReq{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		logkit.Error("Login Req Unmarshal failed", err)
	} else {
		ret := &protocol.LoginRes{}
		if _, nameOk := self.chatHall.UserMapByName.Load(msg.Username); nameOk {
			ret.Result = -1
			ret.Error = "Your name conflicts with another user"
			self.SendMsg(conn, ret)
		} else {
			if self.wordFilter.ContainsAny(msg.Username) {
				ret.Result = -1
				ret.Error = "Your name has some dirty words"
				self.SendMsg(conn, ret)
			} else {
				self.chatHall.Login(msg.Username, msg.Password, connId)
				ret.Result = 0
				self.SendMsg(conn, ret)
			}
		}
	}
}

func (self *RoomMgr) NotifyRoom(room *Room, content string, excludeName string) {
	notify := &protocol.ChatRes{}
	notify.Content = content
	notify.IsSystem = true

	room.UserMapByName.Range(func(key, value interface{}) bool {
		u := value.(*User)
		if u.Conn != nil && u.Name != excludeName {
			self.SendMsg(u.Conn, notify)
		}
		return true
	})
}

func (self *RoomMgr) HandleJoinRoom(data []byte, connId uint32) {
	user := self.chatHall.GetUserByConnId(connId)
	if user == nil {
		return
	}

	msg := &protocol.JoinRoomReq{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		logkit.Error("JoinRoom Req Unmarshal failed", err)
	} else {
		roomId := msg.RoomId
		iRoom, _ := self.roomMap.Load(roomId)
		var room *Room
		if iRoom != nil {
			room = iRoom.(*Room)
		} else {
			room = newRoom(roomId)
		}

		if user.RoomId != 0 {
			self.RemoveUser(user.RoomId, user.Name)
		}

		user.RoomId = roomId
		self.chatHall.UserMapByName.Store(user.Name, user)
		room.JoinUser(user)
		self.NotifyRoom(room, fmt.Sprintf("Welcome user %s Join No.%d Room", user.Name, roomId), user.Name)

		ret := &protocol.JoinRoomRes{}
		ret.ChatList = make([]*protocol.ChatRes, 0, MAX_RECENT_CHAT_MSG)
		for _, chatRes := range room.RecentMsg {
			ret.ChatList = append(ret.ChatList, chatRes)
		}
		self.SendMsg(user.Conn, ret)
	}
}

func resolveTime(seconds int64) (int, int, int, int) {
	day := seconds / 86400
	hour := (seconds - day*86400) / 3600
	min := (seconds - day*86400 - hour*3600) / 60
	sec := seconds - day*86400 - hour*3600 - min*60
	return int(day), int(hour), int(min), int(sec)
}

func (self *RoomMgr) CheckGmCommand(user *User, content string) bool {
	if !user.IsGM {
		return false
	}

	if strings.HasPrefix(content, "/popular") {
		ret := &protocol.GMCommandRes{}
		cmdList := strings.Split(content, " ")
		if len(cmdList) != 2 {
			ret.Result = fmt.Sprintf("Popular GM Command Parament Invalid")
		} else {
			sec, err := strconv.Atoi(cmdList[1])
			if err != nil {
				ret.Result = fmt.Sprintf("Popular GM Command Parament Invalid")
			} else {
				maxStr, maxCount := self.GetTopWord(sec)
				if maxStr == "" || maxCount == 0 {
					ret.Result = fmt.Sprintf("%d Second Has Not Any Max Frequency Word", sec)
				} else {
					ret.Result = fmt.Sprintf("%d Second Send Max Frequency Word : %s ,Count : %d", sec, maxStr, maxCount)
				}
			}
		}

		self.SendMsg(user.Conn, ret)

		return true
	}

	if strings.HasPrefix(content, "/stats") {
		ret := &protocol.GMCommandRes{}
		cmdList := strings.Split(content, " ")
		if len(cmdList) != 2 {
			ret.Result = fmt.Sprintf("Popular GM Command Parament Invalid")
		} else {
			targetName := cmdList[1]
			target := self.chatHall.GetUserByName(targetName)
			if target == nil {
				ret.Result = fmt.Sprintf("User %s is not online", targetName)
			} else {
				nowSec := time.Now().Unix()
				day, hour, min, sec := resolveTime(nowSec - target.LoginTime)
				ret.Result = fmt.Sprintf("User %s onlinetime is %dd %dh %dm %ds", targetName, day, hour, min, sec)
			}
		}

		self.SendMsg(user.Conn, ret)
		return true
	}

	return false
}

func (self *RoomMgr) HandleChat(data []byte, connId uint32) {
	user := self.chatHall.GetUserByConnId(connId)
	if user == nil {
		return
	}

	msg := &protocol.ChatReq{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		logkit.Error("Chat Req Unmarshal failed", err)
	} else {
		if self.CheckGmCommand(user, msg.Content) {
			return
		}

		iRoom, _ := self.roomMap.Load(user.RoomId)
		if iRoom == nil {
			return
		}
		room := iRoom.(*Room)

		notify := &protocol.ChatRes{}
		notify.Content = self.wordFilter.Replace(msg.Content)
		notify.FromName = user.Name

		room.UserMapByName.Range(func(key, value interface{}) bool {
			u := value.(*User)
			if u.Conn != nil {
				self.SendMsg(u.Conn, notify)
			}
			return true
		})

		room.PushMsg(notify)
		self.wordStatChan <- msg.Content
	}
}

func (self *RoomMgr) HandlePrivateChat(data []byte, connId uint32) {
	user := self.chatHall.GetUserByConnId(connId)
	if user == nil {
		return
	}

	msg := &protocol.PrivateChatReq{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		logkit.Error("PrivateChat Req Unmarshal failed", err)
	} else {
		iRoom, _ := self.roomMap.Load(user.RoomId)
		if iRoom == nil {
			return
		}
		room := iRoom.(*Room)
		target := room.GetUserByName(msg.ToName)
		if target == nil {
			ret := &protocol.PrivateChatRes{}
			ret.Result = -1
			ret.Error = fmt.Sprintf("Target %s not in you chat room", msg.ToName)
			self.SendMsg(user.Conn, ret)
			return
		}

		realText := self.wordFilter.Replace(msg.Content)

		if target.Conn != nil {
			chatRes := &protocol.ChatRes{}
			chatRes.FromName = user.Name
			chatRes.Content = realText
			chatRes.IsPrivate = true
			self.SendMsg(target.Conn, chatRes)
		}

		if user.Conn != nil {
			chatRes := &protocol.PrivateChatRes{}
			chatRes.Result = 0
			chatRes.ToName = msg.ToName
			chatRes.Content = realText
			self.SendMsg(user.Conn, chatRes)
		}
	}
}
