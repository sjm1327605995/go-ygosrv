package tcp

import (
	"bytes"
	"github.com/panjf2000/gnet/v2"
	"go-ygosrv/core/duel"
)

type TCPDecoder struct {
}

func (t *TCPDecoder) Decode(buff *bytes.Buffer, player *duel.DuelPlayer) gnet.Action {
	//var (
	//	length    uint16
	//	lengthInt int
	//)
	//err := binary.Read(buff, binary.LittleEndian, &length)
	//if err != nil {
	//	return gnet.None
	//}
	//lengthInt = int(length)
	//if buff.Len() < lengthInt {
	//	arr := buff.Next(lengthInt)
	//
	//}

	duel.HandleCTOSPacket(player, buff.Bytes()[2:])
	buff.Reset()
	return gnet.None
}
