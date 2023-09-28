package duel

import (
	"bytes"
	"go-ygosrv/core/msg/stoc"
	"go-ygosrv/core/ygocore"
	"sync"
)

type SingleDuel struct {
	DuelModeBase
	locker       sync.Locker
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

func (s *SingleDuel) Chat(dp *DuelPlayer, msg BytesMessage) {
	_ = s.Write(dp, stoc.STOC_CHAT, msg)
}

// JoinGame TODO 并发加入房间的问题
func (s *SingleDuel) JoinGame(dp *DuelPlayer, reader *bytes.Buffer) {
	s.locker.Lock()
	defer s.locker.Unlock()
	if s.HostPlayer == nil {
		s.HostPlayer = dp
	}

	s.pplayers[0] = dp
	s.pDuel = ygocore.CreateGame()
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

func (receiver *SingleDuel) UpdateDeck(dp *DuelPlayer, reader *bytes.Buffer) {
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

func (receiver *SingleDuel) Analyze(reader *bytes.Buffer) int {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) Surrender(dp *DuelPlayer) {
	//TODO implement me
	panic("implement me")
}

func (receiver *SingleDuel) GetResponse(dp *DuelPlayer, reader *bytes.Buffer) {
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
