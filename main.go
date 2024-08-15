package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"github.com/sjm1327605995/go-ygosrv/game"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient"
	"sync"
	"time"

	"github.com/sjm1327605995/go-ygosrv/game/ocgcore"
	"github.com/spf13/viper"
	"path/filepath"

	"log"
)

func main() {
	viper.SetDefault("BanlistFile", "lflist.conf")
	viper.SetDefault("ScriptDirectory", "script")
	viper.SetDefault("RootPath", ".")
	viper.SetDefault("DatabaseFile", "cards.cdb")
	game.InitBanListManager(viper.GetString("BanlistFile"))

	_, err := game.InitCardManager(filepath.Join(viper.GetString("RootPath"), viper.GetString("DatabaseFile")))
	if err != nil {
		panic(err)
	}
	ocgcore.InitOcrCore("ocgcore.dll", viper.GetString("ScriptDirectory"), viper.GetString("DatabaseFile"))

	//ClientVersion = Config.GetUInt("ClientVersion", ClientVersion);

	addr := fmt.Sprintf("tcp://127.0.0.1:8080")
	Game := &game.IBaseGame{}
	err = Game.Start()
	if err != nil {
		panic(err)
	}
	//TCP 和UDP 都支持。对TCP分装的。可以通过TCP添加一层协议解析获取内容
	var srv = NewServer()

	log.Println("server exits:", gnet.Run(srv, addr, gnet.WithMulticore(true), gnet.WithReusePort(true), gnet.WithTicker(false)))
}

const (
	TCP uint8 = iota + 1
	WS
)
const poolSize = 10000

type Server struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
	goPool    *ants.Pool
}

type BytesPool struct {
	pool *sync.Pool
}

func NewBytesPool() *BytesPool {
	return &BytesPool{pool: &sync.Pool{
		New: func() any {
			return bytes.NewBuffer(make([]byte, 0, 1024))
		},
	}}
}
func (b *BytesPool) Get() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}
func (b *BytesPool) Put(buffer *bytes.Buffer) {
	buffer.Reset()

	if buffer != nil || buffer.Cap() <= 1024 {
		b.pool.Put(buffer)
	}
	buffer = nil
}
func NewServer() *Server {
	var err error
	goPool, err := ants.NewPool(10000, ants.WithExpiryDuration(time.Second*5))
	if err != nil {
		panic(err)
	}
	return &Server{
		goPool: goPool,
	}
}

func (wss *Server) OnBoot(eng gnet.Engine) gnet.Action {
	wss.eng = eng
	logging.Infof("echo server with multi-core=%t is listening on %s", wss.multicore, wss.addr)
	return gnet.None
}

func (wss *Server) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	ctx := new(Context)
	ctx.msgChan = make(chan *bytes.Reader, 10)
	p := game.NewPlayer(Game, ygoclient.NewClient(c))
	ctx.player = p

	c.SetContext(ctx)
	go ctx.Msg()
	return nil, gnet.None
}

func (wss *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	ctx := c.Context().(*Context)
	ctx.player.Disconnect()
	close(ctx.msgChan)
	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return gnet.None
}

var (
	Game *game.IBaseGame
)

func (wss *Server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	ctx := c.Context().(*Context)
	msgLen, err := c.Next(2)
	if err != nil {
		return gnet.Close
	}
	if len(msgLen) < 2 {
		return gnet.None
	}
	msgLength := int(binary.LittleEndian.Uint16(msgLen))
	if c.InboundBuffered() < msgLength {
		return gnet.None
	}
	arr, err := c.Next(msgLength)
	if err != nil {
		return gnet.Close
	}
	buffer := make([]byte, len(arr))
	copy(buffer, arr)
	reader := bytes.NewReader(buffer)
	ctx.msgChan <- reader
	return gnet.None
}

func (wss *Server) OnTick() (delay time.Duration, action gnet.Action) {
	return 3 * time.Second, gnet.None
}

type tcpReadOp uint8

const (
	readLen tcpReadOp = iota
	readMsg
)

type Context struct {
	player  *game.Player
	msgChan chan *bytes.Reader
}

func (c *Context) Msg() {
	for v := range c.msgChan {
		c.player.Parse(v)
	}
}
