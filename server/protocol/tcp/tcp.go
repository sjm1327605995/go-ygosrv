package tcp

import (
	"bytes"
	"github.com/panjf2000/gnet/v2"
	"go-ygosrv/core/duel"
)

type TCPDecoder struct {
}

func (t *TCPDecoder) Decode(buff *bytes.Buffer, player *duel.DuelPlayer) gnet.Action {
	return gnet.None
}
