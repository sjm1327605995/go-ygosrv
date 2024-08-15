package stoc

import (
	"bytes"
	"encoding/binary"
	"github.com/sjm1327605995/go-ygosrv/msg/host"
)

const ChatMsgLimit = 255 * 2

type ErrorMsg struct {
	Msg       uint8
	_padding1 uint16 //结构体字节填充
	_padding2 uint8  //结构体字节填充
	Code      uint32
}

func (t ErrorMsg) Marshal() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, 100))
	err := binary.Write(b, binary.LittleEndian, &t)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type HandResult struct {
	Res1 uint8
	Res2 uint8
}

func (t HandResult) Marshal() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, 2))
	err := binary.Write(b, binary.LittleEndian, &t)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type CreateGame struct {
	GameId uint32
}

type JoinGame struct {
	Info host.HostInfo
}

func (j JoinGame) Marshal() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, 100))
	err := binary.Write(b, binary.LittleEndian, &j)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type TypeChange struct {
	Type uint8
}

func (t TypeChange) Marshal() ([]byte, error) {
	return []byte{t.Type}, nil
}

//type ExitGame struct {
//	Pos int8
//}
//
//func (p *ExitGame) ToBytes(conn *websocket.Conn)error {
//	return utils.SetData(conn,STOC_E , &p.Pos)
//}

type TimeLimit struct {
	Player   uint8
	LeftTime uint16
}

func (t TimeLimit) Marshal() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, 100))
	err := binary.Write(b, binary.LittleEndian, &t)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type Chat struct {
	Player uint16
	Msg    []byte //256 *2 byte
}

// ToBytes 不使用额外空间原地复制。减少性能浪费
func (c *Chat) ToBytes(buff *bytes.Buffer) error {

	err := binary.Write(buff, binary.LittleEndian, c.Player)
	if err != nil {
		return err
	}
	return binary.Write(buff, binary.LittleEndian, WSStr(c.Msg))

}
func WSStr(arr []byte) []byte {
	i := 0
	for ; i < len(arr); i += 2 {
		if arr[i] == 0 && arr[i+1] == 0 {
			break
		}
	}
	return arr[:i+2]
}

type HSPlayerEnter struct {
	Name [40]byte
	Pos  uint16
}

func (h HSPlayerEnter) Marshal() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, 50))
	err := binary.Write(b, binary.LittleEndian, &h)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type HSWatchChange struct {
	WatchCount uint16
}

func (h HSWatchChange) Marshal() ([]byte, error) {
	arr := make([]byte, 2)
	binary.LittleEndian.PutUint16(arr, h.WatchCount)
	return arr, nil
}

type HSPlayerChange struct {
	Status uint8
}

func (t HSPlayerChange) Marshal() ([]byte, error) {
	return []byte{t.Status}, nil
}
