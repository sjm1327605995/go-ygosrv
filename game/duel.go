package game

import (
	"bytes"
	"io"
	"math/rand"
	"runtime"
)

type Duel struct {
	pDuel        uintptr
	analyzer     func(uint8, *bytes.Reader, []byte) int
	errorHandler func(string)
	buffer       []byte
}

var duels = make(map[uintptr]*Duel)

func (duel *Duel) SetAnalyzer(analyzer func(uint8, *bytes.Reader, []byte) int) {
	duel.analyzer = analyzer
}

func (duel *Duel) SetErrorHandler(errorHandler func(string)) {
	duel.errorHandler = errorHandler
}

func (duel *Duel) InitPlayers(startLP, startHand, drawCount int32) {
	SetPlayerInfo(duel.pDuel, 0, startLP, startHand, drawCount)
	SetPlayerInfo(duel.pDuel, 1, startLP, startHand, drawCount)
}

func (duel *Duel) AddCard(cardID int32, owner int, location byte) {
	NewCard(duel.pDuel, uint32(cardID), byte(owner), byte(owner), location, 0, 0)
}

func (duel *Duel) AddTagCard(cardID int32, owner int, location byte) {
	NewTagCard(duel.pDuel, uint32(cardID), byte(owner), byte(location))
}

func (duel *Duel) Start(options int) {
	StartDuel(duel.pDuel, int32(options))
}

func (duel *Duel) Process() int {
	fail := 0
	for {
		result := Process(duel.pDuel)
		length := int(result & 0xFFFF)
		if length > 0 {
			fail = 0
			duel.buffer = make([]byte, 4096)
			GetMessage(duel.pDuel, duel.buffer)
			res, _ := duel.handleMessage(duel.buffer, length)
			if res != 0 {
				return res
			}
		} else if fail++; fail == 10 {
			return -1
		}
	}
}
func (duel *Duel) QueryFieldInfo() []byte {
	var buf = make([]byte, 256)
	_ = QueryFieldInfo(duel.pDuel, buf)
	//duel.
	//Marshal.Copy(_buffer, result, 0, 256);
	return buf
}

func (duel *Duel) SetResponseInt(resp int) {
	SetResponsei(duel.pDuel, int32(resp))
}

func (duel *Duel) SetResponseByte(resp []byte) {
	if len(resp) > 64 {
		return
	}
	buf := make([]byte, 64)
	copy(buf, resp)
	SetResponseb(duel.pDuel, buf)
}

func (duel *Duel) QueryFieldCount(player int, location byte) int {
	return int(QueryFieldCount(duel.pDuel, byte(player), location))
}

func (duel *Duel) QueryFieldCard(player int, location byte, flag int, useCache bool) []byte {
	var useCache8 int32
	if useCache {
		useCache8 = 1
	}
	length := QueryFieldCard(duel.pDuel, byte(player), location, int32(flag), duel.buffer, useCache8)
	return duel.buffer[:length]
}

func (duel *Duel) QueryCard(player, location, sequence, flag int, useCache bool) []byte {
	var useCache8 int32
	if useCache {
		useCache8 = 1
	}
	length := QueryCard(duel.pDuel, byte(player), byte(location), byte(sequence), int32(flag), duel.buffer, useCache8)
	return duel.buffer[:length]
}

func (duel *Duel) End() {
	EndDuel(duel.pDuel)
	duel.dispose()
}

func (duel *Duel) GetNativePtr() uintptr {
	return duel.pDuel
}

func (duel *Duel) dispose() {
	delete(duels, duel.pDuel)
	duel.buffer = nil
}

func (duel *Duel) handleMessage(raw []byte, len int) (int, error) {
	reader := bytes.NewReader(raw[:len])
	for {
		msg, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return -1, err
		}
		result := -1
		if duel.analyzer != nil {
			result = duel.analyzer(msg, reader, raw)
		}
		if result != 0 {
			return result, nil
		}
	}
	return 0, nil
}

func (duel *Duel) onMessage(messageType uint32) {

	buffer := GetLogMessage(duel.pDuel)
	message := string(bytes.TrimRight(buffer, "\x00"))
	if duel.errorHandler != nil {
		duel.errorHandler(message)
	}
}

func NewDuel(seed uint32) *Duel {
	random := rand.New(rand.NewSource(int64(seed)))
	pDuel := CreateDuel(int32(random.Uint32()))
	return createDuel2(pDuel)
}

func createDuel2(pDuel uintptr) *Duel {
	if pDuel == 0 {
		return nil
	}
	duel := &Duel{
		pDuel:  pDuel,
		buffer: make([]byte, 4096),
	}
	duels[pDuel] = duel
	runtime.SetFinalizer(duel, func(d *Duel) {
		d.dispose()
	})
	return duel
}
