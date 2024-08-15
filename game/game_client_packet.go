package game

import (
	"bytes"
	"encoding/binary"
	"io"
)

type GameClientPacket struct {
	Content []byte
	reader  *bytes.Reader
}

func NewGameClientPacket(content []byte) *GameClientPacket {
	return &GameClientPacket{
		Content: content,
		reader:  bytes.NewReader(content),
	}
}

func (gcp *GameClientPacket) ReadCtos() (uint8, error) {
	var message uint8
	err := binary.Read(gcp.reader, binary.LittleEndian, &message)
	if err != nil {
		return 0, err
	}
	return message, nil
}

func (gcp *GameClientPacket) ReadByte() (byte, error) {
	var b byte
	err := binary.Read(gcp.reader, binary.LittleEndian, &b)
	if err != nil {
		return 0, err
	}
	return b, nil
}

func (gcp *GameClientPacket) ReadToEnd() ([]byte, error) {
	return io.ReadAll(gcp.reader)
}

func (gcp *GameClientPacket) ReadSByte() (int8, error) {
	var b int8
	err := binary.Read(gcp.reader, binary.LittleEndian, &b)
	if err != nil {
		return 0, err
	}
	return b, nil
}

func (gcp *GameClientPacket) ReadInt16() (int16, error) {
	var i int16
	err := binary.Read(gcp.reader, binary.LittleEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (gcp *GameClientPacket) ReadInt32() (int32, error) {
	var i int32
	err := binary.Read(gcp.reader, binary.LittleEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (gcp *GameClientPacket) ReadUInt32() (uint32, error) {
	var i uint32
	err := binary.Read(gcp.reader, binary.LittleEndian, &i)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (gcp *GameClientPacket) ReadUnicode(len int) (string, error) {
	buf := make([]byte, len*2) // Assuming UTF-16 encoding (2 bytes per character)
	_, err := gcp.reader.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (gcp *GameClientPacket) GetPosition() (int64, error) {
	return gcp.reader.Seek(0, io.SeekCurrent)
}

func (gcp *GameClientPacket) SetPosition(pos int64) error {
	_, err := gcp.reader.Seek(pos, io.SeekStart)
	return err
}
