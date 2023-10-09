package protocol

import (
	"bytes"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"go-ygosrv/core/duel"
	"go-ygosrv/server/protocol/tcp"
	"go-ygosrv/server/protocol/websocket"
	"time"
)

type Server struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
}

func (wss *Server) OnBoot(eng gnet.Engine) gnet.Action {
	wss.eng = eng
	logging.Infof("echo server with multi-core=%t is listening on %s", wss.multicore, wss.addr)
	return gnet.None
}

func (wss *Server) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	ctx := new(Context)
	ctx.player = &duel.DuelPlayer{Conn: c}
	c.SetContext(ctx)

	return nil, gnet.None
}

func (wss *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}

	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return gnet.None
}

func (wss *Server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	ws := c.Context().(*Context)
	//先把数据读取到buff中
	if ws.readBufferBytes(c) == gnet.Close {
		return gnet.Close
	}

	switch ws.protocol {
	case 0: //等待解析
		//不包含websocket当做tcp处理
		if bytes.Contains(ws.buf.Bytes(), []byte("Upgrade: websocket")) {
			ws.protocol = duel.WS
			ws.player.Protocol = duel.WS
			ws.Decoder = &websocket.WsDecoder{}
		} else {
			ws.protocol = duel.TCP
			ws.player.Protocol = duel.TCP
			ws.Decoder = &tcp.TCPDecoder{}
		}
	case duel.WS, duel.TCP:

	default:
		return gnet.Close
	}

	return ws.Decoder.Decode(&ws.buf, ws.player)
}

func (wss *Server) OnTick() (delay time.Duration, action gnet.Action) {
	return 3 * time.Second, gnet.None
}

type Decoder interface {
	Decode(buff *bytes.Buffer, player *duel.DuelPlayer) gnet.Action
}
type Context struct {
	protocol uint8
	buf      bytes.Buffer // 从实际socket中读取到的数据缓存
	Decoder  Decoder
	player   *duel.DuelPlayer
}

func (w *Context) readBufferBytes(c gnet.Conn) gnet.Action {
	size := c.InboundBuffered()
	buf := make([]byte, size, size)
	read, err := c.Read(buf)
	if err != nil {
		logging.Infof("read err! %w", err)
		return gnet.Close
	}
	if read < size {
		logging.Infof("read bytes len err! size: %d read: %d", size, read)
		return gnet.Close
	}
	w.buf.Write(buf)
	return gnet.None
}
