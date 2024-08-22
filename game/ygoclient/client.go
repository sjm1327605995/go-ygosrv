package ygoclient

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"io"
)

type YGOClient struct {
	con gnet.Conn
}

func NewClient(con gnet.Conn) *YGOClient {
	return &YGOClient{con: con}
}

type Packet interface {
	Send(writer io.Writer) error
}

func (c *YGOClient) SendBytes(packet []byte) error {
	_, err := c.con.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func (c *YGOClient) Send(packet io.ReadSeeker) error {
	packet.Seek(0, io.SeekStart)
	all, err := io.ReadAll(packet)
	if err != nil {
		return err
	}
	data := make([]byte, len(all)+2)
	binary.LittleEndian.PutUint16(data, uint16(len(all)))
	copy(data[2:], all)
	fmt.Println("hex", hex.EncodeToString(data))
	_, err = c.con.Write(data)
	return err
}

//func (c *YGOClient) SendV(packet Packet, v int) error {
//	err := packet.Send(c.con)
//	if err != nil {
//		return err
//	}
//	_, err = c.con.Write([]byte{byte(v)})
//	return err
//}
