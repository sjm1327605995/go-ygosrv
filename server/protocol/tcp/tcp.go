package tcp

import (
	"bytes"
	"github.com/panjf2000/gnet/v2"
)

type TCPDecoder struct {
}

func (t *TCPDecoder) Decode(c gnet.Conn, buff *bytes.Buffer) gnet.Action {
	return gnet.None
}
