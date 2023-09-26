package duel

import (
	"bytes"
	"encoding/binary"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

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
func HandleCTOSPacket(player *DuelPlayer, data []byte) {
	var (
		buf=bytes.NewBuffer(data)
	)
	proto:=binary.LittleEndian.
}
