package duel

import "sync"

var allDuelRoom sync.Map

type DuelRoom struct {
	locker     sync.RWMutex
	Players    map[string]*DuelPlayer
	HostPlayer *DuelPlayer
	Game       DuelMode
}

func JoinOrCreateRoom(dp *DuelPlayer) *DuelRoom {
	var room = new(DuelRoom)
	val, has := allDuelRoom.LoadOrStore(dp.Pass, room)
	if has {
		room = val.(*DuelRoom)
	}
	room.locker.Lock()
	defer room.locker.Unlock()
	if room.Players == nil {
		room.Players = make(map[string]*DuelPlayer, 2)
		room.HostPlayer = dp
		room.Game = &SingleDuel{}

	}
	dp.Room = room
	oldDp, has := room.Players[dp.Name]
	if has {
		if oldDp == room.HostPlayer {
			oldDp = dp
		}
		_ = oldDp.Conn.Close()
	} else {
		dp.Pos = uint8(len(room.Players))
	}
	room.Players[dp.Name] = dp
	return room
}
func (r *DuelRoom) TypeChange(dp *DuelPlayer) uint8 {
	if r.HostPlayer == dp {
		return 0x10
	}
	return 0
}
func (r *DuelRoom) Broadcast(proto uint8, msg BytesMessage) {

	r.locker.RLock()
	defer r.locker.RUnlock()
	for i := range r.Players {
		_ = r.Game.Write(r.Players[i], proto, msg)
	}
}

func (r *DuelRoom) LeaveGame(dp *DuelPlayer) {

	r.locker.Lock()
	defer r.locker.Unlock()
	delete(r.Players, dp.Name)
	_ = dp.Conn.Close()
}
func (r *DuelRoom) CurrentPlayers() (players []*DuelPlayer) {
	r.locker.RLock()
	defer r.locker.RUnlock()
	for i := range r.Players {
		players = append(players, r.Players[i])
	}
	return
}
