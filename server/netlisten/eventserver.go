package netlisten

import (
	"chatroot/logkit"
	"encoding/binary"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"time"
)

type NetListener struct {
	*gnet.EventServer
	Address    string
	Codec      gnet.ICodec
	WorkerPool *goroutine.Pool
}

func NewNetListener(addr string, port string) *NetListener {
	encoderConfig := gnet.EncoderConfig{
		ByteOrder:                       binary.LittleEndian,
		LengthFieldLength:               2,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}
	decoderConfig := gnet.DecoderConfig{
		ByteOrder:           binary.LittleEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   2,
		LengthAdjustment:    0,
		InitialBytesToStrip: 2,
	}

	codec = gnet.NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)

	listener := &NetListener{
		Address:    fmt.Sprintf("%s:%s", addr, port),
		Codec:      codec,
		WorkerPool: goroutine.Default(),
	}

	return listener
}

func (self *NetListener) Run() error {
	err := gnet.Serve(self, self.Address, gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(self.Codec))
	if err != nil {
		logkit.Panic("Server Start Failed", err.Error())
	}
}

func (self *NetListener) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	logkit.Info("Test codec server is listening on %s (loops: %d)\n", srv.Addr.String(), srv.NumEventLoop)
	return
}

func (self *NetListener) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {

}

func (self *NetListener) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {

}

func (self *NetListener) OnClosed(c gnet.Conn, err error) (action gnet.Action) {

}
