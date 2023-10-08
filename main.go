package main

import (
	"flag"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"go-ygosrv/server/protocol"

	"log"
)

var (
	port      int
	multicore bool
)

func main() {

	// Example command: go run main.go --port 8080 --multicore=true
	flag.IntVar(&port, "port", 8080, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")

	flag.Parse()
	addr := fmt.Sprintf("tcp://127.0.0.1:%d", port)
	//TCP 和UDP 都支持。对TCP分装的。可以通过TCP添加一层协议解析获取内容
	var srv = new(protocol.Server)

	//ygocore.InitCore()
	log.Println("server exits:", gnet.Run(srv, addr, gnet.WithMulticore(multicore), gnet.WithReusePort(true), gnet.WithTicker(false)))
}
