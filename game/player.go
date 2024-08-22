package game

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/network/ctos"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/network/gamestate"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/network/stoc"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/player"
	"github.com/sjm1327605995/go-ygosrv/utils"
	"io"
)

type Player struct {
	Game            *IBaseGame
	Name            string
	IsAuthenticated bool
	Type            int
	Deck            *Deck
	State           int
	client          *ygoclient.YGOClient
}

func NewPlayer(game *IBaseGame, client *ygoclient.YGOClient) *Player {
	return &Player{
		Game:            game,
		Name:            "",
		IsAuthenticated: false,
		Type:            player.Undefined,
		State:           player.None,
		client:          client,
	}

}
func (p *Player) Send(data io.ReadSeeker) error {
	return p.client.Send(data)
}
func (p *Player) Disconnect() {
	//_client.Close();
}
func (p *Player) OnDisconnected() {
	//p.Game.RemovePlayer(this);
}

func (p *Player) SendTypeChange() error {
	packet := NewGamePacketFactory().Create(stoc.TypeChange)
	packet.WriteByte(byte(p.Type + IFElse(p.Game.HostPlayer.Equals(p), player.Host, 0)))
	return p.Send(packet)
}

func (p *Player) Equals(player *Player) bool {
	return p == player
}
func (p *Player) Parse(packet *bytes.Reader) error {
	msgTp, _ := packet.ReadByte()
	var msgErr error
	fmt.Println("msg", msgTp)
	switch msgTp {
	case ctos.PlayerInfo:
		fmt.Println("PlayerInfo")
		msgErr = p.OnPlayerInfo(packet)

	case ctos.JoinGame:
		fmt.Println("JoinGame")
		msgErr = p.OnJoinGame(packet)

	case ctos.CreateGame:
		fmt.Println("CreateGame")
		msgErr = p.OnCreateGame(packet)
	}
	if !p.IsAuthenticated {
		return nil
	}
	switch msgTp {

	case ctos.Chat:
		msgErr = p.OnChat(packet)

	case ctos.HsToDuelist:
		msgErr = p.Game.MoveToDuelist(p)

	case ctos.HsToObserver:
		msgErr = p.Game.MoveToObserver(p)

	case ctos.LeaveGame:
		msgErr = p.Game.RemovePlayer(p)

	case ctos.HsReady:
		msgErr = p.Game.SetReady(p, true)

	case ctos.HsNotReady:
		msgErr = p.Game.SetReady(p, false)

	case ctos.HsKick:
		msgErr = p.OnKick(packet)

	case ctos.HsStart:
		msgErr = p.Game.StartDuel(p)

	case ctos.HandResult:
		msgErr = p.OnHandResult(packet)

	case ctos.TpResult:

		msgErr = p.OnTpResult(packet)

	case ctos.UpdateDeck:

		msgErr = p.OnUpdateDeck(packet)

	case ctos.Response:

		msgErr = p.OnResponse(packet)

	case ctos.Surrender:

		msgErr = p.Game.Surrender(p, 0, false)

	}
	return msgErr
}

type playerInfo struct {
	utils.MessageData
	Name string
}

func (p *Player) OnPlayerInfo(packet *bytes.Reader) error {
	if p.Name != "" {
		return nil
	}
	var err error
	p.Name, err = ReadUnicode(packet, 20)
	return err
}

func (p *Player) OnCreateGame(packet *bytes.Reader) error {
	err := p.Game.SetRules(packet)
	if err != nil {
		return err
	}
	ReadUnicode(packet, 20) //hostname
	ReadUnicode(packet, 30) //password

	err = p.Game.AddPlayer(p)
	if err != nil {
		return err
	}
	p.IsAuthenticated = true
	return nil
}
func (p *Player) OnJoinGame(packet *bytes.Reader) error {

	//if p.Name != "" || p.Type != player.Undefined {
	//	return nil
	//}

	var (
		version   int16
		gameid    int32
		spaceData int16
	)
	err := Read(packet, &version, &gameid, &spaceData)
	if err != nil {
		return err
	}

	err = p.Game.AddPlayer(p)
	if err != nil {
		return err
	}
	p.IsAuthenticated = true
	return nil
}

func (p *Player) OnChat(packet *bytes.Reader) error {
	msg, _ := ReadUnicode(packet, 256)
	return p.Game.Chat(p, msg)

}

func (p *Player) OnKick(packet *bytes.Reader) error {

	pos, err := packet.ReadByte()
	if err != nil {
		return err
	}

	return p.Game.KickPlayer(p, pos)
}

func (p *Player) OnHandResult(packet *bytes.Reader) error {

	res, err := packet.ReadByte()
	if err != nil {
		return err
	}

	return p.Game.HandResult(p, res)
}

func (p *Player) OnTpResult(packet *bytes.Reader) error {
	tp, err := packet.ReadByte()
	if err != nil {
		return err
	}
	return p.Game.TpResult(p, tp != 0)
}

func (p *Player) OnUpdateDeck(packet *bytes.Reader) error {
	if p.Type == player.Observer {
		return nil
	}
	deck := new(Deck)

	var (
		main int32
		side int32
	)
	binary.Read(packet, binary.LittleEndian, &main)
	binary.Read(packet, binary.LittleEndian, &side)
	for i := int32(0); i < main; i++ {
		var cardId int32
		binary.Read(packet, binary.LittleEndian, &cardId)
		deck.AddMain(cardId)
	}
	for i := int32(0); i < side; i++ {
		var cardId int32
		binary.Read(packet, binary.LittleEndian, &cardId)
		deck.AddSide(cardId)
	}
	if p.Game.State == gamestate.Lobby {
		p.Deck = deck
		p.Game.IsReady[p.Type] = false
	} else if gamestate.Side == p.Game.State {
		if p.Game.IsReady[p.Type] {
			return nil
		}
		if !p.Deck.CheckBool(p.Deck) {
			writer := NewGamePacketFactory().Create(stoc.ErrorMsg)
			writer.WriteByte(3)

			writer.WriteByte(0)
			return p.Send(writer)
		}
		p.Deck = deck
		p.Game.IsReady[p.Type] = true
		p.Send(NewGamePacketFactory().Create(stoc.DuelStart))
		p.Game.MatchSide()
	}
	return nil
}

func (p *Player) OnResponse(packet *bytes.Reader) error {
	if p.Game.State != gamestate.Duel {
		return nil
	}

	if p.State != player.Response {
		return nil
	}
	resp, _ := io.ReadAll(packet)

	if len(resp) > 64 {
		return nil
	}
	p.State = player.None
	return p.Game.SetResponseBytes(resp)
}
