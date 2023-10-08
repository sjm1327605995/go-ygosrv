package host

import (
	"bytes"
	"go-ygosrv/utils"
)

const StrLimit = 40

type HostInfo struct {
	Lflist        uint32
	Rule          uint16
	Mode          uint8
	DuleRule      uint8
	NoCheckDeck   uint16
	NoShuffleDeck uint16
	StartLp       uint32
	StartHand     uint8
	DrawCount     uint8
	TimeLimit     uint16
}

func (h *HostInfo) Parse(b *bytes.Buffer) (err error) {
	return utils.GetData(b, &h.Lflist, &h.Rule, &h.Mode, &h.DuleRule, &h.NoCheckDeck, &h.NoShuffleDeck,
		&h.StartLp, &h.StartHand, &h.DrawCount, &h.TimeLimit)

}

type HostPacket struct {
	Identifier uint16
	Version    uint16
	Port       uint16
	IpAddr     uint32
	Name       string //长度为40
	Host       HostInfo
}

func (h *HostPacket) Parse(b *bytes.Buffer) (err error) {

	err = utils.GetData(b, &h.Identifier, &h.Version, &h.Port, &h.IpAddr)
	if err != nil {
		return
	}
	_, h.Name, err = utils.UTF16ToStr(b.Next(StrLimit))
	if err != nil {
		return
	}
	return h.Host.Parse(b)

}

type HostRequest struct {
	Identifier uint16
}

func (h *HostRequest) Parse(b *bytes.Buffer) (err error) {
	return utils.GetData(b, &h.Identifier)

}

type HandResult struct {
	res uint8
}

func (h *HandResult) Parse(b *bytes.Buffer) (err error) {
	return utils.GetData(b, &h.res)

}
