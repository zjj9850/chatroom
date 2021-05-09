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

var msg_chan chan proto.Message

func main() {
	var port = flag.Int("port", 8700, "chat server port")
	var address = flag.String("addr", "0.0.0.0", "chat server address")
	flag.Parse()

	msg_chan = make(chan proto.Message)

	conn := newConnect(*address, *port)
	defer conn.Close()

	go handleinput()
	// go testConn(conn)
	// go connRead(conn)
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

func connRead(conn net.Conn) {
	for {

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

}

func chat(cmdList []string) {

}

func chatTo(cmdList []string) {

}

func handleinput() {
	fmt.Println("Help:")
	fmt.Println("login username userpassord")
	fmt.Println("join roomid")
	fmt.Println("chat content")
	fmt.Println("to toName content")
	fmt.Println("exit")
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
		}
	}
}
