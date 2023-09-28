package tcp

import (
	"encoding/binary"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"go-ygosrv/core/duel"
	"log"
)

type EchoServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	Addr      string
	Multicore bool
}

func (es *EchoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	log.Printf("echo server with multi-core=%t is listening on %s\n", es.Multicore, es.Addr)
	return gnet.None
}

func (es *EchoServer) OnTraffic(c gnet.Conn) gnet.Action {
	player := c.Context().(*duel.DuelPlayer)
	buf, _ := c.Next(-1)
	packetLen := int(binary.LittleEndian.Uint16(buf))
	if packetLen > len(buf)-1 {
		logging.Infof("conn[%v] refuse packet", c.RemoteAddr().String(), buf)
		return gnet.Close
	}
	duel.HandleCTOSPacket(player, buf[2:])
	return gnet.None
}

func (wss *EchoServer) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	player := &duel.DuelPlayer{
		Conn: c,
	}

	c.SetContext(player)

	return nil, gnet.None
}

func (wss *EchoServer) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}
	player := c.Context().(*duel.DuelPlayer)
	duel.Leave(player)
	//TODO 用户离开后离开房间
	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return gnet.None
}
