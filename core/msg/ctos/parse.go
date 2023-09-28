package ctos

import (
	"bytes"
	"go-ygosrv/core/msg/host"
	"go-ygosrv/utils"
)

type PlayerInfo struct {
	Name string
}

const (
	StrLimit = 40
)

func (p *PlayerInfo) Parse(b *bytes.Buffer) (err error) {
	// 将二进制数组转换为字符串
	p.Name, err = utils.UTF16ToStr(b.Next(StrLimit))
	return
}

type TPResult struct {
	Res uint8
}

func (h *TPResult) Parse(b *bytes.Buffer) (err error) {
	return utils.GetData(b, &h.Res)

}

type CreateGame struct {
	Info host.HostInfo
	Name string
	Pass string
}

func (h *CreateGame) Parse(b *bytes.Buffer) (err error) {
	err = h.Info.Parse(b)
	if err != nil {
		return err
	}
	h.Name, err = utils.UTF16ToStr(b.Next(StrLimit))
	if err != nil {
		return
	}
	h.Pass, err = utils.UTF16ToStr(b.Next(StrLimit))
	if err != nil {
		return
	}
	return
}

type JoinGame struct {
	Version uint16
	Align   uint16
	GameId  uint32
	Pass    string
}

// Pass: [40] - 房间密码
func (h *JoinGame) Parse(b *bytes.Buffer) (err error) {
	err = utils.GetData(b, &h.Version, &h.Align, &h.GameId)
	if err != nil {
		return
	}
	h.Pass, err = utils.UTF16ToStr(b.Next(StrLimit))
	if err != nil {
		return
	}
	return
}

type Kick struct {
	Pos uint16
}

func (h *Kick) Parse(b *bytes.Buffer) (err error) {
	return utils.GetData(b, &h.Pos)

}

type HandResult struct {
	Res uint8
}

func (h *HandResult) Parse(b *bytes.Buffer) (err error) {
	return utils.GetData(b, &h.Res)

}
