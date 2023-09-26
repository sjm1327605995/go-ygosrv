package duel

import (
	"fmt"
	"go-ygosrv/core/msg/stoc"
	"go-ygosrv/utils"
)

type SingleDuel struct {
	DuelModeBase
	players      [2]*DuelPlayer
	pplayers     [2]*DuelPlayer
	Ready        [2]bool
	PDeck        [2]Deck
	DeckError    [2]int32
	HandRes      [2]uint8
	LastResponse uint8
	//std::set<DuelPlayer*> observers;
	//	Replay last_replay;
	MatchMode   bool
	MatchKill   int
	DuelCount   uint8
	TpPlayer    uint8
	MatchResult [3]uint8
	TimeLimit   [2]int16
	TimeElapsed int16
}

func (s *SingleDuel) Chat(dp *DuelPlayer, msg []byte) {
	var scc = stoc.Chat{
		Player: dp.Type,
		Msg:    msg,
	}

	err := SendBufferToPlayer(dp, stoc.STOC_CHAT, scc.Player, scc.Msg)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func (receiver *SingleDuel) JoinGame(dp *DuelPlayer, reader *utils.BitReader) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) LeaveGame(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) ToObserver(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) PlayerReady(dp *DuelPlayer, isReady bool) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) PlayerKick(dp *DuelPlayer, pos uint8) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) UpdateDeck(dp *DuelPlayer, reader *utils.BitReader) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) StartDuel(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) HandResult(dp *DuelPlayer, uint82 uint8) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) TPResult(dp *DuelPlayer, uint82 uint8) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) Process() {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) Analyze(reader *utils.BitReader) int {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) Surrender(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) GetResponse(dp *DuelPlayer, reader *utils.BitReader) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) TimeConfirm(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) EndDuel() {
	//TODO implement me
	panic("implement me")
}
