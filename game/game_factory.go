package game

import (
	"encoding/binary"
	"fmt"
	"github.com/sjm1327605995/go-ygosrv/utils"
	"unicode/utf16"
)

type GamePacketFactory struct {
	*utils.MemoryStream
}

func NewGamePacketFactory() *GamePacketFactory {
	return &GamePacketFactory{MemoryStream: utils.NewMemoryStream()}
}

func (gpf *GamePacketFactory) Create(message uint8) *GamePacketFactory {
	_, _ = gpf.MemoryStream.Write([]byte{message})
	return gpf
}
func (gpf *GamePacketFactory) WriteByte(b uint8) *GamePacketFactory {
	_, _ = gpf.MemoryStream.Write([]byte{b})
	return gpf
}

func (gpf *GamePacketFactory) WriteUnicode(text string, length int) {

	// Convert the string to UTF-16 encoded bytes.
	unicodeBytes := utf16.Encode([]rune(text))
	if len(unicodeBytes) > length*2 {
		unicodeBytes = unicodeBytes[:length*2]
		_ = binary.Write(gpf.MemoryStream, binary.LittleEndian, unicodeBytes)
	} else {
		unicodeBytes = append(unicodeBytes, make([]uint16, length-len(unicodeBytes))...)
	}
	_ = binary.Write(gpf.MemoryStream, binary.LittleEndian, unicodeBytes)
}
func (gpf *GamePacketFactory) Write(value any) *GamePacketFactory {
	switch t := value.(type) {
	case string:
		_, _ = gpf.MemoryStream.Write([]byte(t))
	case []byte:
		_, _ = gpf.MemoryStream.Write(t)
	case int:
		fmt.Println("int")
		err := binary.Write(gpf.MemoryStream, binary.LittleEndian, uint8(t))
		if err != nil {
			fmt.Println(err)
		}
	default:
		err := binary.Write(gpf.MemoryStream, binary.LittleEndian, t)
		if err != nil {
			fmt.Println(err)
		}
	}

	return gpf
}
func (gpf *GamePacketFactory) WriteLen(arr []byte, length int) *GamePacketFactory {
	_, _ = gpf.MemoryStream.Write(arr[:length])
	return gpf
}
func (gpf *GamePacketFactory) Bytes() []byte {
	return gpf.MemoryStream.Bytes()
}
