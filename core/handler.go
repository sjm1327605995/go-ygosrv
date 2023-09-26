package core

import (
	"encoding/binary"
	"fmt"
	"go-ygosrv/core/msg/ctos"
	"go-ygosrv/core/msg/stoc"
	"go-ygosrv/utils"

	"github.com/panjf2000/gnet/v2"
)

type Context struct {
	writeFunc func(conn gnet.Conn, arr []byte) error
}

func (c *Context) SetWriter(f func(conn gnet.Conn, arr []byte) error) {
	c.writeFunc = f
}
func (c *Context) Open() {

}
func (c *Context) OnClose() {

}
func (c *Context) HandleMessage(conn gnet.Conn, reader *utils.BitReader) {
	var (
		proto = reader.ReadUint8()
	)
	switch proto {
	case ctos.CTOS_CHAT:
		s, err := utils.UTF16ToStr(reader.Next(40))
		fmt.Println(s, err)
	case ctos.CTOS_JOIN_GAME:
		//	r := room.JoinOrCreate(0, c.Player.Name, conn)
		//	c.Room = r
		err := c.SendPacketToPlayer(conn, stoc.STOC_JOIN_GAME)
		fmt.Println(err)
	case ctos.CTOS_PLAYER_INFO:
		//这里有个很麻烦的问题。原来的ygopro是局域网的游戏对战。所以是根据用户名进行的区分。但是这是服务器所以不能根据用户。只能更具用户名取
		var playerInfo ctos.PlayerInfo
		err := playerInfo.Parse(reader)
		if err != nil {
			fmt.Println(err)
			return
		}
		//c.Player.Name = playerInfo.Name
	}
	return
}
func (c *Context) SendPacketToPlayer(conn gnet.Conn, proto uint8, arr ...byte) error {
	length := uint16(len(arr) + 1)
	var res = make([]byte, length+2)
	binary.LittleEndian.PutUint16(arr, length)
	res[2] = proto
	if length > 3 {
		copy(res[3:], arr)
	}
	return c.writeFunc(conn, res)
}
func (c *Context) SendPacketToPlayer(conn gnet.Conn, proto uint8, arr ...byte) error {
	length := uint16(len(arr) + 1)
	var res = make([]byte, length+2)
	binary.LittleEndian.PutUint16(arr, length)
	res[2] = proto
	if length > 3 {
		copy(res[3:], arr)
	}
	return c.writeFunc(conn, res)
}
