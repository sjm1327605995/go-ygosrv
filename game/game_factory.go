package game

import (
	"encoding/binary"
	"github.com/sjm1327605995/go-ygosrv/utils"

	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/network/stoc"

	"unicode/utf16"
)

type GamePacketFactory struct {
	*utils.Buffer
}

func (gpf *GamePacketFactory) Read(p []byte) (n int, err error) {

	return gpf.Buffer.Read(p)
}

func NewGamePacketFactory() *GamePacketFactory {
	return &GamePacketFactory{Buffer: utils.New(make([]byte, 0))}
}

func (gpf *GamePacketFactory) Create(message uint8) *GamePacketFactory {
	_, _ = gpf.Buffer.Write([]byte{message})
	return gpf
}
func (gpf *GamePacketFactory) WriteByte(b uint8) *GamePacketFactory {
	_, _ = gpf.Buffer.Write([]byte{b})
	return gpf
}
func (gpf *GamePacketFactory) CreateGameMessage(message uint8) *GamePacketFactory {
	_, _ = gpf.Buffer.Write([]byte{stoc.GameMsg, message})
	return gpf
}
func (gpf *GamePacketFactory) WriteUnicode(text string, length int) {
	unicode := utf16.Encode([]rune(text))
	if len(unicode) > length {
		unicode = unicode[:length]
	}
	_ = binary.Write(gpf.Buffer, binary.LittleEndian, unicode)
}
func (gpf *GamePacketFactory) Write(value interface{}) *GamePacketFactory {
	switch t := value.(type) {
	case string:
		_, _ = gpf.Buffer.Write([]byte(t))
	case []byte:
		_, _ = gpf.Buffer.Write(t)
	default:
		_ = binary.Write(gpf.Buffer, binary.LittleEndian, value)
	}

	return gpf
}
func (gpf *GamePacketFactory) WriteLen(arr []byte, length int) *GamePacketFactory {
	_, _ = gpf.Buffer.Write(arr[:length])
	return gpf
}
func (gpf *GamePacketFactory) Bytes() []byte {
	return gpf.Buffer.Bytes()
}
