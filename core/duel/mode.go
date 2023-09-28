package duel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"go-ygosrv/core/msg/stoc"
)

type DuelPlayer struct {
	Name     string //40 byte
	game     DuelMode
	Type     uint16
	Status   uint8
	Protocol uint8
	Conn     gnet.Conn
}

type DuelMode interface {
	Chat(dp *DuelPlayer, msg BytesMessage)
	JoinGame(dp *DuelPlayer, reader *bytes.Buffer)
	LeaveGame(dp *DuelPlayer)
	ToObserver(dp *DuelPlayer)
	PlayerReady(dp *DuelPlayer, isReady bool)
	PlayerKick(dp *DuelPlayer, pos uint8)
	UpdateDeck(dp *DuelPlayer, reader *bytes.Buffer)
	StartDuel(dp *DuelPlayer)
	HandResult(dp *DuelPlayer, uint82 uint8)
	TPResult(dp *DuelPlayer, uint82 uint8)
	Process()
	Analyze(reader *bytes.Buffer) int
	Surrender(dp *DuelPlayer)
	GetResponse(dp *DuelPlayer, reader *bytes.Buffer)
	TimeConfirm(dp *DuelPlayer)
	EndDuel()
	PDuel() int64
	Write(dp *DuelPlayer, proto uint8, msg BytesMessage) error
}
type DuelModeBase struct {
	//Etimer
	HostPlayer *DuelPlayer
	DuelStage  int32
	pDuel      int64
	Name       string //40个字节
	Pass       string //40个字节
}

func (d *DuelModeBase) Chat(dp *DuelPlayer, msg BytesMessage) {
	_ = d.Write(dp, stoc.STOC_CHAT, msg)
}

func (d *DuelModeBase) JoinGame(dp *DuelPlayer, reader *bytes.Buffer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) LeaveGame(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) ToObserver(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) PlayerReady(dp *DuelPlayer, isReady bool) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) PlayerKick(dp *DuelPlayer, pos uint8) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) UpdateDeck(dp *DuelPlayer, reader *bytes.Buffer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) StartDuel(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) HandResult(dp *DuelPlayer, uint82 uint8) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) TPResult(dp *DuelPlayer, uint82 uint8) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) Process() {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) Analyze(reader *bytes.Buffer) int {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) Surrender(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) GetResponse(dp *DuelPlayer, reader *bytes.Buffer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) TimeConfirm(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) EndDuel() {
	//TODO implement me
	panic("implement me")
}
func (d *DuelModeBase) PDuel() int64 {
	return 0
}

type ParseMessage interface {
	Parse(*bytes.Buffer) error
}
type BytesMessage interface {
	ToBytes(*bytes.Buffer) error
}

func (d *DuelModeBase) Write(dp *DuelPlayer, proto uint8, msg BytesMessage) error {
	buffer := bytes.NewBuffer(make([]byte, 3, 100))
	err := msg.ToBytes(buffer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	arr := buffer.Bytes()
	binary.LittleEndian.PutUint16(arr, uint16(len(arr)-2))
	arr[2] = proto
	switch dp.Protocol {
	//Websocket
	case 0:
		return wsutil.WriteServerMessage(dp.Conn, ws.OpBinary, arr)
	//TCP
	case 1:
		_, err := dp.Conn.Write(arr)
		return err
	}
	return nil
}
