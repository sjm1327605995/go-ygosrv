package duel

import (
	"bytes"
	"fmt"
	"go-ygosrv/core/msg/ctos"
	"go-ygosrv/core/msg/stoc"
)

var model DuelModeBase

//HandleCTOSPacket 重构dp结构体优化调用链
func HandleCTOSPacket(dp *DuelPlayer, data []byte) {
	var (
		buf = bytes.NewBuffer(data[1:])
	)

	pktType := data[0]
	if (pktType != ctos.CTOS_SURRENDER) && (pktType != ctos.CTOS_CHAT) && (dp.Status == 0xff || (dp.Status == 1 && dp.Status != pktType)) {
		return
	}

	switch pktType {
	case ctos.CTOS_RESPONSE:
		if dp.game == nil || dp.game.PDuel() == 0 {
			return
		}
		dp.game.GetResponse(dp, buf)
	case ctos.CTOS_TIME_CONFIRM:
		if dp.game == nil || dp.game.PDuel() == 0 {
			return
		}
		dp.game.TimeConfirm(dp)
	case ctos.CTOS_CHAT:
		//客户端发过来的消息
		var chat = stoc.Chat{
			Player: dp.Type,
			Msg:    buf.Bytes(),
		}
		if dp.game != nil {
			//开始了游戏
			dp.game.Chat(dp, &chat)
		}
	case ctos.CTOS_UPDATE_DECK:
		if dp.game == nil {
			return
		}

		dp.game.UpdateDeck(dp, buf)
	case ctos.CTOS_HAND_RESULT:
		if dp.game == nil {
			return
		}
		var tpRes ctos.HandResult
		err := tpRes.Parse(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		dp.game.HandResult(dp, tpRes.Res)
	case ctos.CTOS_TP_RESULT:
		if dp.game == nil {
			return
		}
		var pkt ctos.TPResult
		err := pkt.Parse(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		dp.game.TPResult(dp, pkt.Res)
	case ctos.CTOS_PLAYER_INFO:
		var pkt ctos.PlayerInfo
		err := pkt.Parse(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		dp.RealName = pkt.RealName
		dp.Name = pkt.Name
	case ctos.CTOS_CREATE_GAME: //TODO 暂时请求未使用到 比较疑惑
		if dp.game != nil {
			return
		}

	case ctos.CTOS_JOIN_GAME: //TODO 现在如果game为空就进行初始化

		var joinGame ctos.JoinGame
		err := joinGame.Parse(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		room := JoinOrCreate(0, dp)
		if room == nil {
			fmt.Println("room create fail")
			return
		}
		fmt.Println(joinGame)
		if dp.game == nil {
			d := &SingleDuel{}
			d.players[0] = dp
			dp.game = d

		}
		//var join = stoc.JoinGame{Info: host.HostInfo{
		//	Lflist:        0,
		//	Rule:          0,
		//	Mode:          0,
		//	DuleRule:      0,
		//	NoCheckDeck:   false,
		//	NoShuffleDeck: false,
		//	StartLp:       0,
		//	StartHand:     0,
		//	DrawCount:     0,
		//	TimeLimit:     0,
		//}}
		//暂时不知道什么意思 web客户端未使用到
		dp.game.Write(dp, stoc.STOC_JOIN_GAME, &BytesMsg{67, 63, 66, 112, 0, 0, 5, 0, 0, 0, 0, 0, 64, 31, 0, 0, 5, 1, 240, 0})
		dp.game.Write(dp, stoc.STOC_TYPE_CHANGE, &BytesMsg{16})
		//C++ 和C中都是以0为结尾。为了兼容C所以做的字符串末尾标识
		s := BytesMsg(append(WSStr(dp.RealName), 0, 0))
		dp.game.Write(dp, stoc.STOC_HS_PLAYER_ENTER, &s)
	case ctos.CTOS_LEAVE_GAME:
		dp.game.LeaveGame(dp)
		Leave(dp)
	}
}
func WSStr(arr []byte) []byte {
	var i int
	for ; i < len(arr)-1; i = i + 2 {
		if arr[i] == 0 && arr[i+1] == 0 {
			break
		}
	}
	return arr[:i]
}
