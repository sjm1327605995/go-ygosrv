package room

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
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
	clients   map[string]gnet.Conn
}

func JoinOrCreate(roomId uint64, playerId string, conn gnet.Conn) *Room {

	//如果用户存在对应的room直接加入
	//  不存在 如果指定房间id不为0 就加入指定房间 否则创建新的房间号
	//  依照生成的房间号查找是否存在，不存在就创建。存在直接获取到
	rid, has := playerRoom.Load(playerId)
	if has {
		roomId = rid.(uint64)
	} else if roomId == 0 {
		roomId, _ = utils.Sf.NextID()
	}
	playerRoom.Store(playerId, roomId)
	var r = &Room{Id: roomId}
	val, has := roomAll.LoadOrStore(roomId, r)
	if has {
		r = val.(*Room)
	} else {
		r.clients = make(map[string]gnet.Conn, 2)
	}
	r.rwLocker.Lock()
	defer r.rwLocker.Unlock()
	oldConn, has := r.clients[playerId]
	if has {
		err := oldConn.Close()
		if err != nil {
			fmt.Println(err)
		}

	}
	r.clients[playerId] = conn
	fmt.Println("player:", playerId, "join the room:", roomId)
	return r
}
func (r *Room) IsOwner(playerId string) bool {
	return playerId == r.owner
}
func (r *Room) Leave(playerId string) {
	r.rwLocker.Lock()
	defer r.rwLocker.Unlock()
	delete(r.clients, playerId)
	if len(r.clients) == 0 {
		roomAll.Delete(r.Id)
	}
}
