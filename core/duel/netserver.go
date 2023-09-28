package duel

import (
	"bytes"
	"fmt"
	"go-ygosrv/core/msg/ctos"
	"go-ygosrv/core/msg/stoc"
)

var model DuelModeBase

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
		dp.Name = pkt.Name
	case ctos.CTOS_CREATE_GAME: //TODO 暂时请求未使用到 比较疑惑
		if dp.game != nil {
			return
		}

	case ctos.CTOS_JOIN_GAME: //TODO 现在如果game为空就进行初始化
		if dp.game != nil {
			return
		}
		var joinGame ctos.JoinGame
		err := joinGame.Parse(buf)
		if err != nil {
			fmt.Println(err)
			return
		}
		dp.game = &SingleDuel{}
		dp.game.JoinGame(dp, buf)
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
