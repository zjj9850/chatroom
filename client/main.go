package main

import "flag"
import "net"
import "fmt"
import "os"
import "strings"
import "chatclient/protocol"
import "bufio"
import "google.golang.org/protobuf/proto"
import "encoding/binary"
import "time"
import "strconv"

var msg_chan chan proto.Message

var mapHandler map[string]func([]byte)

func main() {
	var port = flag.Int("port", 8700, "chat server port")
	var address = flag.String("addr", "0.0.0.0", "chat server address")
	flag.Parse()

	msg_chan = make(chan proto.Message)
	mapHandler = make(map[string]func([]byte))
	mapHandler[string((&protocol.LoginRes{}).ProtoReflect().Descriptor().FullName())] = handleLoginRes
	mapHandler[string((&protocol.JoinRoomRes{}).ProtoReflect().Descriptor().FullName())] = handleJoinRes
	mapHandler[string((&protocol.PrivateChatRes{}).ProtoReflect().Descriptor().FullName())] = handlePrivateRes
	mapHandler[string((&protocol.ChatRes{}).ProtoReflect().Descriptor().FullName())] = handleChatRes
	mapHandler[string((&protocol.GMCommandRes{}).ProtoReflect().Descriptor().FullName())] = handleGmRes

	conn := newConnect(*address, *port)
	defer conn.Close()

	go handleinput()
	// go testConn(conn)
	go connRead(conn)
	go connWrite(conn)

	select {}
}

func testConn(conn net.Conn) {
	req := &protocol.LoginReq{}
	req.Username = "ABC"
	req.Password = "ABC"
	send_bytes := pack_message(req)

	writer := bufio.NewWriterSize(conn, 1024)
	for {
		l := uint16(len(send_bytes))
		writer.Flush()
		binary.Write(writer, binary.LittleEndian, l)
		writer.Write(send_bytes)
		time.Sleep(time.Second)
	}
}

func newConnect(addr string, port int) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return conn
}

func handleLoginRes(data []byte) {
	res := &protocol.LoginRes{}
	proto.Unmarshal(data, res)
	if res.Result == 0 {
		fmt.Println("Login Success")
	} else {
		fmt.Println("Login Failed:", res.Error)
	}
}

func handleJoinRes(data []byte) {
	res := &protocol.JoinRoomRes{}
	proto.Unmarshal(data, res)
	for _, chat := range res.ChatList {
		fmt.Printf("User:%s Speek:%s\n", chat.FromName, chat.Content)
	}
}

func handlePrivateRes(data []byte) {
	res := &protocol.PrivateChatRes{}
	proto.Unmarshal(data, res)
	if res.Result == 0 {
		fmt.Printf("You Speek To %s : %s\n", res.ToName, res.Content)
	} else {
		fmt.Println("Private Msg Failed", res.Error)
	}
}

func handleChatRes(data []byte) {
	res := &protocol.ChatRes{}
	proto.Unmarshal(data, res)
	if res.IsSystem {
		fmt.Printf("System:%s\n", res.Content)
	} else if res.IsPrivate {
		fmt.Printf("User %s Speek To You : %s\n", res.FromName, res.Content)
	} else {
		fmt.Printf("User %s Speek : %s\n", res.FromName, res.Content)
	}
}

func handleGmRes(data []byte) {
	res := &protocol.GMCommandRes{}
	proto.Unmarshal(data, res)
	fmt.Printf("GMCommand Result:%s\n", res.Result)
}

func connRead(conn net.Conn) {
	for {
		reply := make([]byte, 2048)
		_, err := conn.Read(reply)

		if err != nil {
			fmt.Println("Read Failed", err.Error())
			os.Exit(1)
		}

		l_buf := reply[0:2]
		l := binary.LittleEndian.Uint16(l_buf)
		if l != 0 {
			msg_bytes := reply[2 : 2+l]
			msg := &protocol.Message{}
			err := proto.Unmarshal(msg_bytes, msg)
			if err != nil {
				fmt.Println("Unmarshal failed", string(msg_bytes))
			} else {
				if handler, e := mapHandler[msg.Type]; e {
					handler(msg.Data)
				}
			}
		}
	}
}

func connWrite(conn net.Conn) {
	writer := bufio.NewWriterSize(conn, 1024)
	for {
		select {
		case msg := <-msg_chan:
			send_bytes := pack_message(msg)
			l := uint16(len(send_bytes))
			writer.Flush()
			binary.Write(writer, binary.LittleEndian, l)
			writer.Write(send_bytes)
			writer.Flush()
		}
	}
}

func pack_message(msg proto.Message) []byte {
	packed := &protocol.Message{}
	packed.Type = string(msg.ProtoReflect().Descriptor().FullName())
	packed.Data, _ = proto.Marshal(msg)
	ret, err := proto.Marshal(packed)
	if err != nil {
		return nil
	}
	return ret
}

func login(cmdList []string) {
	if len(cmdList) != 2 {
		return
	}

	req := &protocol.LoginReq{}
	req.Username = cmdList[0]
	req.Password = cmdList[1]
	msg_chan <- req
}

func join(cmdList []string) {
	if len(cmdList) != 1 {
		return
	}

	roomId, err := strconv.Atoi(cmdList[0])
	if err != nil {
		return
	}

	req := &protocol.JoinRoomReq{}
	req.RoomId = uint32(roomId)
	msg_chan <- req
}

func gm(cmdList []string) {
	gm := ""
	for _, cmd := range cmdList {
		gm += cmd + " "
	}
	req := &protocol.ChatReq{}
	req.Content = gm
	msg_chan <- req
}

func chat(cmdList []string) {
	if len(cmdList) < 1 {
		return
	}

	cmd := ""
	for _, msg := range cmdList {
		cmd += msg + " "
	}
	cmd = strings.TrimSuffix(cmd, " ")
	if len(cmd) == 0 {
		return
	}

	req := &protocol.ChatReq{}
	req.Content = cmd
	msg_chan <- req
}

func chatTo(cmdList []string) {
	if len(cmdList) < 2 {
		return
	}

	cmd := ""
	for _, msg := range cmdList[1:] {
		cmd += msg + " "
	}
	cmd = strings.TrimSuffix(cmd, " ")
	if len(cmd) == 0 {
		return
	}

	req := &protocol.PrivateChatReq{}
	req.ToName = cmdList[0]
	req.Content = cmd
	msg_chan <- req
}

func handleinput() {
	fmt.Println(".......................")
	fmt.Println("Help:")
	fmt.Println("Login Cmd: 		login username userpassord")
	fmt.Println("Join Room Cmd:		join roomid")
	fmt.Println("Chat Cmd:			chat content")
	fmt.Println("PrivateChat Cmd:	to toName content")
	fmt.Println("GM Cmd:			gm /stats username")
	fmt.Println("GM Cmd:			gm /popular second")
	fmt.Println("Exit:				exit")
	fmt.Println(".......................")
	fmt.Println("")
	inputReader := bufio.NewReader(os.Stdin)

	for {
		input, err := inputReader.ReadString('\n')
		if err != nil {
			break
		}

		text := strings.TrimSuffix(input, "\n")
		if text == "exit" {
			os.Exit(0)
		}

		textList := strings.Split(text, " ")
		if len(textList) <= 1 {
			continue
		}

		cmd := textList[0]
		switch cmd {
		case "login":
			login(textList[1:])
		case "join":
			join(textList[1:])
		case "chat":
			chat(textList[1:])
		case "to":
			chatTo(textList[1:])
		case "gm":
			gm(textList[1:])
		}
	}
}
