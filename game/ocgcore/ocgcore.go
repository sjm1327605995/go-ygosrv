package ocgcore

import "C"
import (
	"fmt"
	"github.com/sjm1327605995/go-ygocore"
	"github.com/sjm1327605995/go-ygosrv/game"
	"os"
	"path/filepath"
)

var YGOCore *ygocore.YGOCore

func Init(libPath string) {

}
func CreateDuel(seed int32) uintptr {
	return YGOCore.CreateDuel(seed)

}
func StartDuel(pduel uintptr, options int32) {
	YGOCore.StartDuel(pduel, options)
}
func EndDuel(pduel uintptr) {
	YGOCore.EndDuel(pduel)
}
func SetPlayerInfo(pduel uintptr, playerId, lp, startCount, drawCount int32) {
	YGOCore.SetPlayerInfo(pduel, playerId, lp, startCount, drawCount)
}

const (
	LogMessageBufLen = 1024
	MessageBufLen    = 0x1000
	QueryCardBufLen  = 0x2000
	ResponsebBufLen  = 64
)

// GetLogMessage 返回[]byte 长度固定为1024
func GetLogMessage(pduel uintptr) []byte {
	var buf = make([]byte, LogMessageBufLen)
	YGOCore.GetLogMessage(pduel, buf)
	return buf
}

func GetMessage(pduel uintptr, buff []byte) int32 {
	return YGOCore.GetMessage(pduel, buff)

}
func Process(pduel uintptr) int32 {
	return YGOCore.Process(pduel)
}
func NewCard(pduel uintptr, code uint32, owner, playerid, location, sequence, position uint8) {
	YGOCore.NewCard(pduel, code, owner, playerid, location, sequence, position)
}
func NewTagCard(pduel uintptr, code uint32, owner, position uint8) {
	//YGOCore.New
	//C.new_tag_card(C.longlong(pduel), C.uint32_t(code), C.uint8_t(owner), C.uint8_t(position))
}

// QueryCard  buf 长度要大于 0x2000
func QueryCard(pduel uintptr, playerid, location, sequence uint8, queryFlag int32, buf []byte, useCache int32) int32 {
	return YGOCore.QueryCard(pduel, playerid, location, sequence, queryFlag, buf, useCache)
}

func QueryFieldCount(pduel uintptr, playerid, location uint8) int32 {
	return YGOCore.QueryFieldCount(pduel, playerid, location)
}
func QueryFieldCard(pduel uintptr, playerid, location uint8, queryFlag int32, buf []byte, useCache int32) int32 {
	return YGOCore.QueryFieldCard(pduel, playerid, location, queryFlag, buf, useCache)
}

func QueryFieldInfo(pduel uintptr, buf []byte) int32 {
	return YGOCore.QueryFieldInfo(pduel, buf)
}
func SetResponsei(pduel uintptr, value int32) {
	YGOCore.SetResponseI(pduel, value)
}
func SetResponseb(pduel uintptr, buf []byte) {
	YGOCore.SetResponseB(pduel, buf)
}
func PreloadScript(pduel uintptr, script string) int32 {
	return YGOCore.PreloadScript(pduel, script, int32(len(script)))
}

var (
	scriptPath       string
	dataBaseFilePath string
)

func InitOcrCore(libPath, scriptPath string, databaseFile string) {

	YGOCore = ygocore.NewYGOCore(libPath,
		func(scriptName string) (data []byte) {
			base := filepath.Base(scriptName)
			data, _ = os.ReadFile(filepath.Join(scriptPath, base))
			return
		}, func(cardId int32) *ygocore.CardData {
			card := game.GlobalCardManager.GetCard(cardId)
			if card == nil {
				return nil
			}
			return card.CardData
		}, func(msg string) {
			fmt.Printf("message{%s}\n", msg)
		})

}
