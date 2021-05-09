package main

import (
	"bufio"
	"chatroom/filter"
	"chatroom/logkit"
	"chatroom/netlisten"
	"chatroom/room"
	"flag"
	"io"
	"os"
)

const (
	DIRTY_WORDS_FILE = "config/dirtywords.list"
)

func initWordFilter(wordFilter *filter.WordFilter) error {
	file, err := os.Open(DIRTY_WORDS_FILE)
	if err != nil {
		return err
	}
	defer file.Close()

	br := bufio.NewReader(file)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		wordFilter.InsertKeyWords(string(line))
	}

	return nil
}

func main() {
	var port = flag.Int("port", 8700, "chat server port")
	flag.Parse()

	wordFilter := filter.NewWordFilter()
	if err := initWordFilter(wordFilter); err != nil {
		logkit.Panic("Word Filter Init Failed", err)
	} else {
		logkit.Info("Word Filter Init Success")
	}

	server := netlisten.NewNetListener("tcp://0.0.0.0", *port)

	err := make(chan error)

	roomMgr := room.NewRoomMgr(server, wordFilter)
	roomMgr.Init()
	roomMgr.Run()

	err <- server.Run()

	select {
	case e := <-err:
		logkit.Debug("Chat Server is shutdown", e.Error())
	}
}
