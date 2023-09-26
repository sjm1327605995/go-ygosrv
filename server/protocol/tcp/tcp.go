package tcp

import (
	"github.com/panjf2000/gnet/v2"
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
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}
