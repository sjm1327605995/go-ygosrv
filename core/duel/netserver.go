package duel

import (
	"encoding/binary"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func SendBufferToPlayer(dp *DuelPlayer, extraDataLen int, proto uint8, arr ...byte) error {
	var res = make([]byte, extraDataLen+3)
	binary.LittleEndian.PutUint16(arr, uint16(extraDataLen+1))
	res[2] = proto
	if extraDataLen > 0 {
		copy(res[3:], arr)
	}
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
