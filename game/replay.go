package game

import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
)

const (
	FlagCompressed = 0x1
	FlagTag        = 0x2
	MaxReplaySize  = 0x20000
)

type Replay struct {
	Disabled bool
	Header   ReplayHeader
	Writer   *bytes.Buffer

	memoryStream *bytes.Buffer
	data         []byte
}

func NewReplay(seed uint32, tag bool) *Replay {
	replay := &Replay{
		Header: ReplayHeader{
			Id: 0x31707279,
			//	Version: uint32(Program.Config.ClientVersion),
			Flag: 0,
			Seed: seed,
		},
		memoryStream: new(bytes.Buffer),
		Writer:       new(bytes.Buffer),
	}

	if tag {
		replay.Header.Flag |= FlagTag
	}
	return replay
}

func (r *Replay) Check() {
	if r.memoryStream.Len() >= MaxReplaySize {
		r.Writer = nil
		r.memoryStream = nil
		r.Disabled = true
	}
}
func (r *Replay) WriteUnicode(text string, length int) {
	unicode := utf16.Encode([]rune(text))
	if len(unicode) > length {
		unicode = unicode[:length]
	}
	_ = binary.Write(r.Writer, binary.LittleEndian, unicode)
}
func (r *Replay) End() {
	if r.Disabled {
		return
	}

	raw := r.memoryStream.Bytes()

	r.Header.DataSize = uint32(len(raw))
	r.Header.Flag |= FlagCompressed
	var props [8]byte

	// Assuming a function encodeLZMA that performs LZMA encoding
	compressed, _ := encodeLZMA(raw, props)

	writer := new(bytes.Buffer)

	binary.Write(writer, binary.LittleEndian, r.Header.Id)
	binary.Write(writer, binary.LittleEndian, r.Header.Version)
	binary.Write(writer, binary.LittleEndian, r.Header.Flag)
	binary.Write(writer, binary.LittleEndian, r.Header.Seed)
	binary.Write(writer, binary.LittleEndian, r.Header.DataSize)
	binary.Write(writer, binary.LittleEndian, r.Header.Hash)

	writer.Write(props[:])

	writer.Write(compressed)

	r.data = writer.Bytes()
}

func (r *Replay) GetFile() []byte {
	return r.data
}

// Placeholder for the LZMA encoding function
func encodeLZMA(data []byte, props [8]byte) ([]byte, error) {
	// LZMA encoding should be implemented here
	return data, nil
}
