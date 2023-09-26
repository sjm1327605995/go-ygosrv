package host

import (
	"go-ygosrv/utils"
)

const StrLimit = 40

type HostInfo struct {
	Lflist        uint16
	Rule          uint8
	Mode          uint8
	DuleRule      uint8
	NoCheckDeck   bool
	NoShuffleDeck bool
	StartLp       uint32
	StartHand     uint8
	DrawCount     uint8
	TimeLimit     uint16
}

func (h *HostInfo) Parse(b *utils.BitReader) (err error) {
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

func (h *HostPacket) Parse(b *utils.BitReader) (err error) {

	err = utils.GetData(b, &h.Identifier, &h.Version, &h.Port, &h.IpAddr)
	if err != nil {
		return
	}
	h.Name, err = utils.UTF16ToStr(b.Next(StrLimit))
	if err != nil {
		return
	}
	return h.Host.Parse(b)

}

type HostRequest struct {
	Identifier uint16
}

func (h *HostRequest) Parse(b *utils.BitReader) (err error) {
	return utils.GetData(b, &h.Identifier)

}

type HandResult struct {
	res uint8
}

func (h *HandResult) Parse(b *utils.BitReader) (err error) {
	return utils.GetData(b, &h.res)

}
