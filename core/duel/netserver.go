package duel

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"go-ygosrv/core/msg/ctos"
	"go-ygosrv/core/msg/stoc"
	"go-ygosrv/utils"
)

var model DuelModeBase

func SendBufferToPlayer(dp *DuelPlayer, proto uint8, data ...interface{}) error {
	b := bytes.NewBuffer(make([]byte, 2, 50))
	err := binary.Write(b, binary.LittleEndian, proto)
	if err != nil {
		return err
	}
	for i := range data {
		err = binary.Write(b, binary.LittleEndian, data[i])
		if err != nil {
			return err
		}
	}
	arr := b.Bytes()
	binary.LittleEndian.PutUint16(arr, uint16(b.Len()-2))
	return dp.game.Write(dp, arr)
}

func HandleCTOSPacket(dp *DuelPlayer, data []byte) {
	var (
		buf = utils.NewBitReader(data[1:], len(data)-1)
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
		arr := WSStr(buf.Next(256 * 2))
		buff := bytes.NewBuffer(make([]byte, 0, uint16(len(arr)+3+2)))
		err := binary.Write(buff, binary.LittleEndian, uint16(4+len(arr)+1))
		if err != nil {
			fmt.Println(err)
		}
		err = binary.Write(buff, binary.LittleEndian, stoc.STOC_CHAT)
		if err != nil {
			fmt.Println(err)
		}
		var chat = stoc.Chat{
			Player: dp.Type,
			Msg:    arr,
		}

		_ = binary.Write(buff, binary.LittleEndian, &chat)
		model.Chat(dp, buff)
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
