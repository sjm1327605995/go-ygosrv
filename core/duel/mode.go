package duel

import (
	"bytes"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/panjf2000/gnet/v2"
	"go-ygosrv/utils"
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
	Chat(dp *DuelPlayer, reader *bytes.Buffer)
	JoinGame(dp *DuelPlayer, reader *utils.BitReader)
	LeaveGame(dp *DuelPlayer)
	ToObserver(dp *DuelPlayer)
	PlayerReady(dp *DuelPlayer, isReady bool)
	PlayerKick(dp *DuelPlayer, pos uint8)
	UpdateDeck(dp *DuelPlayer, reader *utils.BitReader)
	StartDuel(dp *DuelPlayer)
	HandResult(dp *DuelPlayer, uint82 uint8)
	TPResult(dp *DuelPlayer, uint82 uint8)
	Process()
	Analyze(reader *utils.BitReader) int
	Surrender(dp *DuelPlayer)
	GetResponse(dp *DuelPlayer, reader *utils.BitReader)
	TimeConfirm(dp *DuelPlayer)
	EndDuel()
	PDuel() int64
	Write(dp *DuelPlayer, arr []byte) error
}
type DuelModeBase struct {
	//Etimer
	HostPlayer *DuelPlayer
	DuelStage  int32
	pDuel      int64
	Name       string //40个字节
	Pass       string //40个字节
}

func (d *DuelModeBase) Chat(dp *DuelPlayer, reader *bytes.Buffer) {
	_ = d.Write(dp, reader.Bytes())
}

func (d *DuelModeBase) JoinGame(dp *DuelPlayer, reader *utils.BitReader) {
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

func (d *DuelModeBase) UpdateDeck(dp *DuelPlayer, reader *utils.BitReader) {
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

func (d *DuelModeBase) Analyze(reader *utils.BitReader) int {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) Surrender(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (d *DuelModeBase) GetResponse(dp *DuelPlayer, reader *utils.BitReader) {
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
func (d *DuelModeBase) Write(dp *DuelPlayer, arr []byte) error {
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
