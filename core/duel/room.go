package duel

import (
	"fmt"
	"go-ygosrv/utils"
	"sync"
)

var (
	roomAll    sync.Map
	playerRoom sync.Map
)

type Room struct {
	Id        uint64
	owner     string //playerName 代表房主
	DuelStage int
	pDuel     int64
	rwLocker  sync.RWMutex
	clients   map[string]*DuelPlayer
}

func JoinOrCreate(roomId uint64, player *DuelPlayer) *Room {

	//如果用户存在对应的room直接加入
	//  不存在 如果指定房间id不为0 就加入指定房间 否则创建新的房间号
	//  依照生成的房间号查找是否存在，不存在就创建。存在直接获取到
	rid, has := playerRoom.Load(player.Name)
	if has {
		roomId = rid.(uint64)
	} else if roomId == 0 {
		roomId, _ = utils.Sf.NextID()
	}
	playerRoom.Store(player.Name, roomId)
	var r = &Room{Id: roomId}
	val, has := roomAll.LoadOrStore(roomId, r)
	if has {
		r = val.(*Room)
	} else {
		r.clients = make(map[string]*DuelPlayer, 2)
	}
	r.rwLocker.Lock()
	defer r.rwLocker.Unlock()
	oldPlayer, has := r.clients[player.Name]
	if has {
		err := oldPlayer.Conn.Close()
		if err != nil {
			fmt.Println(err)
		}
		oldPlayer.Conn = player.Conn
		player = oldPlayer

	}
	r.clients[player.Name] = player
	fmt.Println("player:", player.Name, "join the room:", roomId)
	return r
}
func (r *Room) IsOwner(playerId string) bool {
	return playerId == r.owner
}
func Leave(player *DuelPlayer) {
	rid, has := playerRoom.Load(player.Name)
	if !has {
		return
	}
	val, has := playerRoom.Load(rid)
	if !has {
		return
	}
	delete(val.(*Room).clients, player.Name)
	_ = player.Conn.Close()
}
