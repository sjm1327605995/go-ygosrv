package duel

import (
	"bytes"
	"fmt"
	"go-ygosrv/core/msg/ctos"
	"go-ygosrv/core/msg/host"
	"go-ygosrv/core/msg/stoc"
)

const (
	WS uint8 = iota + 1
	TCP
)

var model DuelModeBase

// HandleCTOSPacket 重构dp结构体优化调用链
func HandleCTOSPacket(dp *DuelPlayer, data []byte, length uint16) {
	var (
		buf = bytes.NewBuffer(data[1:])
	)

	pktType := data[0]
	if (pktType != ctos.CTOS_SURRENDER) && (pktType != ctos.CTOS_CHAT) && (dp.Status == 0xff || (dp.Status == 1 && dp.Status != pktType)) {
		return
	}

	switch pktType {
	//case ctos.CTOS_RESPONSE:
	//	if dp.Room.Game == nil || dp.game.PDuel() == 0 {
	//		return
	//	}
	//	dp.Room.GetResponse(dp, buf)
	//case ctos.CTOS_TIME_CONFIRM:
	//	if dp.Room == nil || dp.Room.PDuel() == 0 {
	//		return
	//	}
	//	dp.Room.TimeConfirm(dp)
	//case ctos.CTOS_CHAT:
	//	//客户端发过来的消息
	//	var chat = stoc.Chat{
	//		Player: dp.Type,
	//		Msg:    buf.Bytes(),
	//	}
	//	if dp.game != nil {
	//		//开始了游戏
	//		dp.game.Chat(dp, &chat)
	//	}
	case ctos.CTOS_UPDATE_DECK:
		if dp.Room == nil || dp.Room.Game == nil {
			return
		}
		//
		err := dp.Room.Game.UpdateDeck(dp, buf, length)
		if err != nil {
			fmt.Println("cards err")

		}
	//case ctos.CTOS_HAND_RESULT:
	//	if dp.game == nil {
	//		return
	//	}
	//	var tpRes ctos.HandResult
	//	err := tpRes.Parse(buf)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	dp.game.HandResult(dp, tpRes.Res)
	//case ctos.CTOS_TP_RESULT:
	//	if dp.game == nil {
	//		return
	//	}
	//	var pkt ctos.TPResult
	//	err := pkt.Parse(buf)
	//	if err != nil {
	//		fmt.Println(err)
	//		return
	//	}
	//	dp.game.TPResult(dp, pkt.Res)
	case ctos.CTOS_PLAYER_INFO:
		var pkt ctos.PlayerInfo
		err := pkt.Parse(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		dp.RealName = pkt.RealName
		dp.Name = pkt.Name
	//case ctos.CTOS_CREATE_GAME: //TODO 暂时请求未使用到 比较疑惑
	//	if dp.game != nil {
	//		return
	//	}

	case ctos.CTOS_JOIN_GAME: //TODO 现在如果game为空就进行初始化

		var (
			joinGame ctos.JoinGame
		)
		err := joinGame.Parse(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		dp.Pass = joinGame.Pass
		var (
			duelRoom   = JoinOrCreateRoom(dp)
			typeChange = stoc.TypeChange{Type: duelRoom.TypeChange(dp)}
		)

		duelRoom.Game.Write(dp, stoc.STOC_JOIN_GAME, &stoc.JoinGame{Info: host.HostInfo{
			Lflist:        1883389763,
			Rule:          0,
			Mode:          0,
			DuleRule:      5,
			NoCheckDeck:   0,
			NoShuffleDeck: 0,
			StartLp:       8000,
			StartHand:     5,
			DrawCount:     1,
			TimeLimit:     240,
		}})
		//暂不考虑观战者
		var scpe stoc.HSPlayerEnter

		//不是第一个进入房间的玩家
		if dp.Pos != 0 {
			players := duelRoom.CurrentPlayers()
			for i := range players {
				scpe.Name = players[i].RealName
				scpe.Pos = uint16(players[i].Pos)
				duelRoom.Game.Write(dp, stoc.STOC_HS_PLAYER_ENTER, &scpe)
			}
		}
		scpe.Name = dp.RealName
		scpe.Pos = uint16(dp.Pos)
		duelRoom.Game.Write(dp, stoc.STOC_TYPE_CHANGE, &typeChange)
		duelRoom.Broadcast(stoc.STOC_HS_PLAYER_ENTER, &scpe)
	case ctos.CTOS_LEAVE_GAME:
		dp.Room.LeaveGame(dp)

	}
}
func WSStr(arr []byte) []byte {
	var i int
	for ; i < len(arr)-1; i = i + 2 {
		if arr[i] == 0 && arr[i+1] == 0 {
			break
		}
	}
	for {
		if i > len(arr)-1 {
			break
		}
		arr[i] = 0
		i++
	}
	return arr
}
