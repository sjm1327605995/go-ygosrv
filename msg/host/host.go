package host

import (
	"bytes"
	"encoding/binary"
)

const StrLimit = 40

type HostInfo struct {
	Lflist        uint32
	Rule          uint8
	Mode          uint8
	DuleRule      uint8
	NoCheckDeck   bool
	NoShuffleDeck bool
	Unknown       uint16
	Unknown1      uint8
	StartLp       int32
	StartHand     uint8
	DrawCount     uint8
	TimeLimit     uint16
}

func (h *HostInfo) Parse(b []byte) (err error) {
	reader := bytes.NewReader(b)
	return binary.Read(reader, binary.LittleEndian, h)

}

type HostPacket struct {
	Identifier uint16
	Version    uint16
	Port       uint16
	IpAddr     uint32
	Name       []byte //长度为40
	Host       HostInfo
}

func (h *HostPacket) Parse(buff []byte) (err error) {
	h.Name = make([]byte, StrLimit)
	reader := bytes.NewReader(buff)
	return binary.Read(reader, binary.LittleEndian, h)

}

type HostRequest struct {
	Identifier uint16
}

func (h *HostRequest) Parse(buff []byte) (err error) {
	reader := bytes.NewReader(buff)
	return binary.Read(reader, binary.LittleEndian, h)

}

type HandResult struct {
	res uint8
}

func (h *HandResult) Parse(buff []byte) (err error) {
	reader := bytes.NewReader(buff)
	return binary.Read(reader, binary.LittleEndian, h)
}
