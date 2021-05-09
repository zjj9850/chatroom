package netlisten

import (
	"chatroom/logkit"
	"encoding/binary"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/pool/goroutine"
	"sync/atomic"
	"time"
)

const MAX_FRAME_ChANNEL = 1024

type NetFrame struct {
	Frame  []byte
	ConnId uint32
}

type NetListener struct {
	*gnet.EventServer
	Address      string
	Codec        gnet.ICodec
	WorkerPool   *goroutine.Pool
	ConnID       uint32
	FrameChannel chan *NetFrame
	OpenChannel  chan gnet.Conn
	CloseChannel chan uint32
}

func NewNetListener(addr string, port int) *NetListener {
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

	codec := gnet.NewLengthFieldBasedFrameCodec(encoderConfig, decoderConfig)

	listener := &NetListener{
		Address:    fmt.Sprintf("%s:%d", addr, port),
		Codec:      codec,
		WorkerPool: goroutine.Default(),
		ConnID:     0,
	}

	listener.FrameChannel = make(chan *NetFrame, MAX_FRAME_ChANNEL)
	listener.OpenChannel = make(chan gnet.Conn)
	listener.CloseChannel = make(chan uint32)

	return listener
}

func (self *NetListener) Run() error {
	err := gnet.Serve(self, self.Address, gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*5), gnet.WithCodec(self.Codec))
	if err != nil {
		return err
	}
	return nil
}

func (self *NetListener) SendMsg(c gnet.Conn, msg []byte) {
	_ = self.WorkerPool.Submit(func() {
		c.AsyncWrite(msg)
	})
}

func (self *NetListener) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	logkit.Infof("Chatroom server is listening on %s (loops: %d)", srv.Addr.String(), srv.NumEventLoop)
	return
}

func (self *NetListener) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	netFrame := &NetFrame{
		ConnId: c.Context().(uint32),
	}
	netFrame.Frame = make([]byte, len(frame))
	copy(netFrame.Frame, frame)
	self.FrameChannel <- netFrame
	return
}

func (self *NetListener) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	connId := atomic.AddUint32(&self.ConnID, 1)
	c.SetContext(connId)
	self.OpenChannel <- c
	logkit.Infof("Socket %s has been opened,ConnId:%d", c.RemoteAddr().String(), c.Context().(uint32))
	return
}

func (self *NetListener) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	connId := c.Context().(uint32)
	self.CloseChannel <- connId
	logkit.Infof("Socket %s is closing,ConnId:%d", c.RemoteAddr().String(), connId)
	return
}
