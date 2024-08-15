package game

import (
	"bytes"
	"encoding/binary"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/card"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/msg"
	"io"
)

type GameAnalyser struct {
	Game *IBaseGame
}

func (g *GameAnalyser) Analyser(msgTp uint8, reader *bytes.Reader, raw []byte) int {

	switch msgTp {
	case msg.Retry:
		g.OnRetry()
		return 1
	case msg.Hint:
		g.OnHint(msgTp, reader)
	case msg.Win:
		g.OnWin(msgTp, reader)
		return 2
	case msg.SelectBattleCmd:
		g.OnSelectBattleCmd(msgTp, reader)
		return 1
	case msg.SelectEffectYn:
		g.OnSelectEffectYn(msgTp, reader)
		return 1
	case msg.SelectYesNo:
		g.OnSelectYesNo(msgTp, reader)
		return 1
	case msg.SelectOption:
		g.OnSelectOption(msgTp, reader)
		return 1
	case msg.SelectCard:
		fallthrough
	case msg.SelectTribute:
		g.OnSelectCard(msgTp, reader)
		return 1
	case msg.SelectUnselect:
		g.OnSelectUnselect(msgTp, reader)
		return 1
	case msg.SelectChain:
		return g.OnSelectChain(msgTp, reader)
	case msg.SelectPlace:
		fallthrough
	case msg.SelectDisfield:
		fallthrough
	case msg.SelectPosition:
		g.OnSelectPlace(msgTp, reader)
		return 1
	case msg.SelectCounter:
		g.OnSelectCounter(msgTp, reader)
		return 1
	case msg.SelectSum:
		g.OnSelectSum(msgTp, reader)
		return 1
	case msg.SortCard:
		fallthrough
	case msg.SortChain:
		g.OnSortCard(msgTp, reader)
		return 1
	case msg.ConfirmDecktop:
		g.OnConfirmDecktop(msgTp, reader)

	case msg.ConfirmExtratop:
		g.OnConfirmExtratop(msgTp, reader)

	case msg.ConfirmCards:
		g.OnConfirmCards(msgTp, reader)

	case msg.ShuffleDeck:
		fallthrough
	case msg.RefreshDeck:

		g.SendToAllLen(msgTp, reader, 1)
		fallthrough
	case msg.ShuffleHand:
		g.OnShuffleHand(msgTp, reader)

	case msg.ShuffleExtra:
		g.OnShuffleExtra(msgTp, reader)
	case msg.SwapGraveDeck:
		g.OnSwapGraveDeck(msgTp, reader)

	case msg.ReverseDeck:
		g.SendToAllLen(msgTp, reader, 0)

	case msg.DeckTop:
		g.SendToAllLen(msgTp, reader, 6)

	case msg.ShuffleSetCard:
		g.OnShuffleSetCard(msgTp, reader)

	case msg.NewTurn:
		g.OnNewTurn(msgTp, reader)

	case msg.NewPhase:
		g.OnNewPhase(msgTp, reader)

	case msg.Move:
		g.OnMove(msgTp, reader)

	case msg.PosChange:
		g.OnPosChange(msgTp, reader)

	case msg.Set:
		g.OnSet(msgTp, reader)

	case msg.Swap:
		g.SendToAllLen(msgTp, reader, 16)

	case msg.FieldDisabled:
		g.SendToAllLen(msgTp, reader, 4)

	case msg.Summoned:
		fallthrough
	case msg.SpSummoned:
		fallthrough
	case msg.FlipSummoned:
		g.SendToAllLen(msgTp, reader, 0)
		g.Game.RefreshMonsters(0, nil)
		g.Game.RefreshMonsters(1, nil)
		g.Game.RefreshSpells(0, nil)
		g.Game.RefreshSpells(1, nil)

	case msg.Summoning:
	case msg.SpSummoning:
		g.SendToAllLen(msgTp, reader, 8)

	case msg.FlipSummoning:
		g.OnFlipSummoning(msgTp, reader)

	case msg.Chaining:
		g.SendToAllLen(msgTp, reader, 16)

	case msg.Chained:
		g.SendToAllLen(msgTp, reader, 1)
		g.Game.RefreshAll()

	case msg.ChainSolving:
		g.SendToAllLen(msgTp, reader, 1)

	case msg.ChainSolved:
		g.SendToAllLen(msgTp, reader, 1)
		g.Game.RefreshAll()

	case msg.ChainEnd:
		g.SendToAllLen(msgTp, reader, 0)
		g.Game.RefreshAll()

	case msg.ChainNegated:
		fallthrough
	case msg.ChainDisabled:
		g.SendToAllLen(msgTp, reader, 1)

	case msg.CardSelected:
		g.OnCardSelected(msgTp, reader)

	case msg.RandomSelected:
		g.OnRandomSelected(msgTp, reader)

	case msg.BecomeTarget:
		g.OnBecomeTarget(msgTp, reader)

	case msg.Draw:
		g.OnDraw(msgTp, reader)

	case msg.Damage:
		fallthrough
	case msg.Recover:
		fallthrough
	case msg.LpUpdate:
		fallthrough
	case msg.PayLpCost:
		g.OnLpUpdate(msgTp, reader)

	case msg.Equip:
		g.SendToAllLen(msgTp, reader, 8)

	case msg.Unequip:
		g.SendToAllLen(msgTp, reader, 4)

	case msg.CardTarget:
		fallthrough
	case msg.CancelTarget:
		g.SendToAllLen(msgTp, reader, 8)

	case msg.AddCounter:
		fallthrough
	case msg.RemoveCounter:
		g.SendToAllLen(msgTp, reader, 7)

	case msg.Attack:
		g.SendToAllLen(msgTp, reader, 8)

	case msg.Battle:
		g.SendToAllLen(msgTp, reader, 26)

	case msg.AttackDisabled:
		g.SendToAllLen(msgTp, reader, 0)

	case msg.DamageStepStart:
	case msg.DamageStepEnd:
		g.SendToAllLen(msgTp, reader, 0)
		g.Game.RefreshMonsters(0, nil)
		g.Game.RefreshMonsters(1, nil)

	case msg.MissedEffect:
		g.OnMissedEffect(msgTp, reader)

	case msg.TossCoin:
		fallthrough
	case msg.TossDice:
		g.OnTossCoin(msgTp, reader)

	case msg.RockPaperScissors:
		g.OnRockPaperScissors(msgTp, reader)
		return 1
	case msg.HandResult:
		g.SendToAllLen(msgTp, reader, 1)

	case msg.AnnounceRace:
		g.OnAnnounceRace(msgTp, reader)
		return 1
	case msg.AnnounceAttrib:
		g.OnAnnounceAttrib(msgTp, reader)
		return 1
	case msg.AnnounceCard:
		g.OnAnnounceCard(msgTp, reader)
		return 1
	case msg.AnnounceNumber:
		g.OnAnnounceNumber(msgTp, reader)
		return 1
	case msg.AnnounceCardFilter:
		g.OnAnnounceCardFilter(msgTp, reader)
		return 1
	case msg.CardHint:
		g.SendToAllLen(msgTp, reader, 9)

	case msg.PlayerHint:
		g.SendToAllLen(msgTp, reader, 6)

	case msg.MatchKill:
		g.OnMatchKill(msgTp, reader)

	case msg.TagSwap:
		g.OnTagSwap(msgTp, reader)

	default:
		//throw new Exceptig.On("[GameAnalyser] Unhandled packet id: " + msg);
	}
	return 0
}

func (g *GameAnalyser) OnRetry() {
	player := g.Game.WaitForResponse()
	g.Game.CurPlayers[player].Send(NewGamePacketFactory().Create(msg.Retry))
	g.Game.Replay.End()
}

func (g *GameAnalyser) OnHint(msgTp uint8, reader *bytes.Reader) {
	tp, _ := reader.ReadByte()
	player, _ := reader.ReadByte()
	_, _ = reader.Seek(4, io.SeekCurrent)
	//
	//reader.ReadInt32()

	packet := NewGamePacketFactory().Create(msgTp)
	_, _ = io.Copy(packet.Buffer, reader)

	switch tp {
	case 1, 2, 3, 4, 5:
		g.Game.CurPlayers[player].Send(packet)

	case 6, 7, 8, 9:
		g.Game.SendToAllButInt(packet, int(player))

	case 10:
		if g.Game.IsTag {
			g.Game.CurPlayers[player].Send(packet)
		} else {
			g.Game.SendToAll(packet)
		}
	}
}
func (g *GameAnalyser) OnWin(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	reason, _ := reader.ReadByte()
	g.Game.MatchSaveResult(int(player), int(reason))
	g.SendToAll(msgTp, reader)
}

func (g *GameAnalyser) OnSelectBattleCmd(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	count, _ := reader.ReadByte()
	reader.Seek(int64(count*11), io.SeekCurrent)
	//reader.ReadBytes(count * 11)
	count, _ = reader.ReadByte()
	reader.Seek(int64(count*8+2), io.SeekCurrent)
	//reader.ReadBytes(count*8 + 2)
	g.Game.RefreshAll()
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnSelectIdleCmd(msgTp uint8, reader *bytes.Reader) {

	player, _ := reader.ReadByte()

	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*7, io.SeekCurrent)
	//msg.Reader.ReadBytes(count * 7)
	count, _ = reader.ReadByte()
	//msg.Reader.ReadBytes(count * 7)
	reader.Seek(int64(count)*7, io.SeekCurrent)
	count, _ = reader.ReadByte()
	reader.Seek(int64(count)*7, io.SeekCurrent)
	count, _ = reader.ReadByte()
	reader.Seek(int64(count)*7, io.SeekCurrent)
	count, _ = reader.ReadByte()
	reader.Seek(int64(count)*7, io.SeekCurrent)
	count, _ = reader.ReadByte()
	reader.Seek(int64(count)*11+3, io.SeekCurrent)
	//msg.Reader.ReadBytes(count*11 + 3)

	g.Game.RefreshAll()
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnSelectEffectYn(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	//msg.Reader.ReadBytes(13)
	reader.Seek(13, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnSelectYesNo(msgTp uint8, reader *bytes.Reader) {

	player, _ := reader.ReadByte()
	reader.Seek(4, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnSelectOption(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*4, io.SeekCurrent)
	//msg.Reader.ReadBytes(count * 4)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnSelectCard(msgTp uint8, reader *bytes.Reader) {

	packet := NewGamePacketFactory().Create(msgTp)

	player, _ := reader.ReadByte()
	packet.Write(player)
	io.CopyN(packet.Buffer, reader, 3)

	count, _ := reader.ReadByte()
	packet.Write(count)

	var i uint8
	for i = 0; i < count; i++ {
		var (
			code              int32
			pl, loc, seq, pos byte
		)
		binary.Read(reader, binary.LittleEndian, &code)
		binary.Read(reader, binary.LittleEndian, &pl)
		binary.Read(reader, binary.LittleEndian, &loc)
		binary.Read(reader, binary.LittleEndian, &seq)
		binary.Read(reader, binary.LittleEndian, &pos)

		packet.Write(IFElse(pl == player, code, 0))
		packet.Write(pl)
		packet.Write(loc)
		packet.Write(seq)
		packet.Write(pos)
	}

	g.Game.WaitForResponseN(int(player))
	g.Game.CurPlayers[player].Send(packet)
}

func (g *GameAnalyser) OnSelectUnselect(msgTp uint8, reader *bytes.Reader) {

	packet := NewGamePacketFactory().Create(msgTp)

	player, _ := reader.ReadByte()
	packet.Write(player)
	io.CopyN(packet.Buffer, reader, 4)

	count, _ := reader.ReadByte()

	packet.Write(count)

	var i uint8
	for i = 0; i < count; i++ {
		var (
			code              int32
			pl, loc, seq, pos byte
		)
		binary.Read(reader, binary.LittleEndian, &code)
		binary.Read(reader, binary.LittleEndian, &pl)
		binary.Read(reader, binary.LittleEndian, &loc)
		binary.Read(reader, binary.LittleEndian, &seq)
		binary.Read(reader, binary.LittleEndian, &pos)

		packet.Write(IFElse(pl == player, code, 0))
		packet.Write(pl)
		packet.Write(loc)
		packet.Write(seq)
		packet.Write(pos)
	}

	g.Game.WaitForResponseN(int(player))
	g.Game.CurPlayers[player].Send(packet)
}

func (g *GameAnalyser) OnSelectChain(msgTp uint8, reader *bytes.Reader) int {
	player, _ := reader.ReadByte()
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*13+10, io.SeekCurrent)

	if count > 0 {
		g.Game.WaitForResponseN(int(player))
		g.SendToPlayer(msgTp, reader, int(player))
		return 1
	}

	g.Game.SetResponse(-1)
	return 0
}

func (g *GameAnalyser) OnSelectPlace(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	reader.Seek(5, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnSelectCounter(msgTp uint8, reader *bytes.Reader) {

	player, _ := reader.ReadByte()
	reader.Seek(4, io.SeekCurrent)
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*9, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))

}

func (g *GameAnalyser) OnSelectSum(msgTp uint8, reader *bytes.Reader) {
	reader.Seek(1, io.SeekCurrent)

	player, _ := reader.ReadByte()
	reader.Seek(6, io.SeekCurrent)
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*11, io.SeekCurrent)
	count, _ = reader.ReadByte()
	reader.Seek(int64(count)*11, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnSortCard(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()

	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*7, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnConfirmDecktop(msgTp uint8, reader *bytes.Reader) {
	reader.Seek(1, io.SeekCurrent)
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*7, io.SeekCurrent)
	g.SendToAll(msgTp, reader)
}

func (g *GameAnalyser) OnConfirmExtratop(msgTp uint8, reader *bytes.Reader) {
	reader.Seek(1, io.SeekCurrent)
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*7, io.SeekCurrent)
	g.SendToAll(msgTp, reader)
}

func (g *GameAnalyser) OnConfirmCards(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()

	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*7, io.SeekCurrent)

	packet := NewGamePacketFactory().Create(msgTp)
	buffer, _ := io.ReadAll(reader)
	io.Copy(packet.Buffer, reader)

	if buffer[7] == card.Hand {
		g.Game.SendToAll(packet)
	} else {
		g.Game.CurPlayers[player].Send(packet)
	}
}

func (g *GameAnalyser) OnShuffleHand(msgTp uint8, reader *bytes.Reader) {
	packet := NewGamePacketFactory().Create(msgTp)
	player, _ := reader.ReadByte()

	count, _ := reader.ReadByte()
	packet.Buffer.Write([]byte{player, count})
	reader.Seek(int64(count)*4, io.SeekCurrent)
	var i byte
	for i = 0; i < count; i++ {
		packet.Write(0)
	}

	g.SendToPlayer(msgTp, reader, int(player))
	g.Game.SendToAllButInt(packet, int(player))
	g.Game.RefreshHand(int(player), nil)
}
func (g *GameAnalyser) OnShuffleExtra(msgTp uint8, reader *bytes.Reader) {
	packet := NewGamePacketFactory().Create(msgTp)
	player, _ := reader.ReadByte()

	count, _ := reader.ReadByte()
	packet.Buffer.Write([]byte{player, count})

	reader.Seek(int64(count)*4, io.SeekCurrent)
	var i byte
	for i = 0; i < count; i++ {
		packet.Write(0)
	}
	g.SendToPlayer(msgTp, reader, int(player))
	g.Game.SendToAllButInt(packet, int(player))
	g.Game.RefreshHand(int(player), nil)
}

func (g *GameAnalyser) OnSwapGraveDeck(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	g.SendToAll(msgTp, reader)
	g.Game.RefreshGrave(int(player), nil)
}

func (g *GameAnalyser) OnShuffleSetCard(msgTp uint8, reader *bytes.Reader) {

	location, _ := reader.ReadByte()

	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*8, io.SeekCurrent)

	g.SendToAll(msgTp, reader)
	if location == card.MonsterZone {
		g.Game.RefreshMonsters(0, nil)
		g.Game.RefreshMonsters(1, nil)
	} else {
		g.Game.RefreshSpells(0, nil)
		g.Game.RefreshSpells(1, nil)
	}
}

func (g *GameAnalyser) OnNewTurn(msgTp uint8, reader *bytes.Reader) {
	g.Game.TimeReset()
	if !g.Game.IsTag {
		g.Game.RefreshAll()
	}

	player, _ := reader.ReadByte()
	g.Game.CurrentPlayer = int(player)
	g.SendToAll(msgTp, reader)

	g.Game.TurnCount++
}

func (g *GameAnalyser) OnNewPhase(msgTp uint8, reader *bytes.Reader) {
	reader.Seek(2, io.SeekCurrent)
	g.SendToAll(msgTp, reader)
	g.Game.RefreshAll()
}

func (g *GameAnalyser) OnMove(msgTp uint8, reader *bytes.Reader) {
	var raw = make([]byte, 16)

	_, _ = reader.Read(raw)

	pc := int(raw[4])
	pl := int(raw[5])
	cc := int(raw[8])
	cl := int(raw[9])
	cs := int(raw[10])
	cp := int(raw[11])

	g.SendToPlayer(msgTp, reader, cc)

	packet := NewGamePacketFactory().Create(msgTp)
	packet.Write(raw)
	if !((cl&card.Grave|card.Overlay) != 0) && ((cl&card.Deck|card.Hand) != 0) || (cp&card.FaceDown) != 0 {
		packet.Seek(2, io.SeekStart)
		packet.Write(0)
	}
	g.Game.SendToAllButInt(packet, cc)

	if cl != 0 && (cl&0x80) == 0 && (cl != pl || pc != cc) {
		g.Game.RefreshSingle(cc, cl, cs)
	}

}

func (g *GameAnalyser) OnPosChange(msgTp uint8, reader *bytes.Reader) {
	var raw = make([]byte, 9)

	_, _ = reader.Read(raw)

	g.SendToAll(msgTp, reader)

	var (
		cc = raw[4]

		cl = raw[5]

		cs = raw[6]

		pp = raw[7]

		cp = raw[8]
	)

	if (pp&card.FaceDown != 0) && (cp&card.FaceUp) != 0 {
		g.Game.RefreshSingle(int(cc), int(cl), int(cs))
	}

}

func (g *GameAnalyser) OnSet(msgTp uint8, reader *bytes.Reader) {
	reader.Seek(4, io.SeekCurrent)
	var raw = make([]byte, 4)
	_, _ = reader.Read(raw)

	packet := NewGamePacketFactory().Create(msg.Set)
	packet.Write(0)
	packet.Write(raw)
	g.Game.SendToAll(packet)
}

func (g *GameAnalyser) OnFlipSummoning(msgTp uint8, reader *bytes.Reader) {
	var raw = make([]byte, 8)
	_, _ = reader.Read(raw)
	g.Game.RefreshSingle(int(raw[4]), int(raw[5]), int(raw[6]))
	g.SendToAll(msgTp, reader)
}

func (g *GameAnalyser) OnCardSelected(msgTp uint8, reader *bytes.Reader) {
	_, _ = reader.ReadByte()
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*4, io.SeekCurrent)
}

func (g *GameAnalyser) OnRandomSelected(msgTp uint8, reader *bytes.Reader) {
	_, _ = reader.ReadByte()
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*4, io.SeekCurrent)
	g.SendToAll(msgTp, reader)
}

func (g *GameAnalyser) OnBecomeTarget(msgTp uint8, reader *bytes.Reader) {
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*4, io.SeekCurrent)
	g.SendToAll(msgTp, reader)
}

func (g *GameAnalyser) OnDraw(msgTp uint8, reader *bytes.Reader) {

	packet := NewGamePacketFactory().Create(msgTp)

	player, _ := reader.ReadByte()

	count, _ := reader.ReadByte()
	packet.Buffer.Write([]byte{player, count})
	var i byte
	for i = 0; i < count; i++ {
		var (
			code uint32
		)
		_ = binary.Read(reader, binary.LittleEndian, &code)

		if (code & 0x80000000) != 0 {
			packet.Write(code)
		} else {
			packet.Write(0)
		}
	}
	g.SendToTeam(msgTp, reader, int(player))
	g.SendToOpponentTeam(packet, int(player))
	g.Game.SendToObservers(packet)
}

func (g *GameAnalyser) OnLpUpdate(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	var (
		value int32
	)
	binary.Read(reader, binary.LittleEndian, &value)

	switch msgTp {
	case msg.LpUpdate:
		g.Game.LifePoints[player] = value
	case msg.PayLpCost:
		fallthrough
	case msg.Damage:
		g.Game.LifePoints[player] -= value
		if g.Game.LifePoints[player] < 0 {
			g.Game.LifePoints[player] = 0
		}

	case msg.Recover:
		g.Game.LifePoints[player] += value

	}

	g.SendToAll(msgTp, reader)
}

func (g *GameAnalyser) OnMissedEffect(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	reader.Seek(7, io.SeekCurrent)

	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnTossCoin(msgTp uint8, reader *bytes.Reader) {
	reader.ReadByte()

	count, _ := reader.ReadByte()
	reader.Seek(int64(count), io.SeekCurrent)

	g.SendToAll(msgTp, reader)
}

func (g *GameAnalyser) OnRockPaperScissors(msgTp uint8, reader *bytes.Reader) {

	player, _ := reader.ReadByte()
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnAnnounceRace(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	reader.Seek(5, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnAnnounceAttrib(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	reader.Seek(5, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnAnnounceCard(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	reader.Seek(4, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnAnnounceNumber(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*4, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnAnnounceCardFilter(msgTp uint8, reader *bytes.Reader) {
	player, _ := reader.ReadByte()
	count, _ := reader.ReadByte()
	reader.Seek(int64(count)*4, io.SeekCurrent)
	g.Game.WaitForResponseN(int(player))
	g.SendToPlayer(msgTp, reader, int(player))
}

func (g *GameAnalyser) OnMatchKill(msgTp uint8, reader *bytes.Reader) {
	reader.Seek(4, io.SeekCurrent)
	if g.Game.IsMatch {
		g.Game.MatchKill()
		g.SendToAll(msgTp, reader)
	}
}

func (g *GameAnalyser) OnTagSwap(msgTp uint8, reader *bytes.Reader) {

	packet := NewGamePacketFactory().Create(msg.TagSwap)

	player, _ := reader.ReadByte()
	packet.WriteByte(player)
	io.CopyN(packet.Buffer, reader, 1) // mcount

	ecount, _ := reader.ReadByte()
	packet.WriteByte(ecount)
	pcount, _ := reader.ReadByte()
	packet.WriteByte(pcount) // pcount

	hcount, _ := reader.ReadByte()
	packet.WriteByte(hcount)

	io.CopyN(packet.Buffer, reader, 4) // topcode
	n := int(hcount) + int(ecount)
	for i := 0; i < n; i++ {
		var code uint32
		binary.Read(reader, binary.LittleEndian, &code)
		if (code & 0x80000000) != 0 {
			packet.Write(code)
		} else {
			packet.Write(0)
		}
	}

	if g.Game.CurPlayers[player].Equals(g.Game.Players[player*2]) {
		g.Game.CurPlayers[player] = g.Game.Players[player*2+1]
	} else {
		g.Game.CurPlayers[player] = g.Game.Players[player*2]
	}
	g.SendToPlayer(msgTp, reader, int(player))
	g.Game.SendToAllButInt(packet, int(player))

	g.Game.RefreshExtra(int(player), nil)
	g.Game.RefreshMonsters(0, nil)
	g.Game.RefreshMonsters(1, nil)
	g.Game.RefreshSpells(0, nil)
	g.Game.RefreshSpells(1, nil)
	g.Game.RefreshHand(0, nil)
	g.Game.RefreshHand(1, nil)
}

func (g *GameAnalyser) SendToAll(msg uint8, reader *bytes.Reader) {

	packet := NewGamePacketFactory().Create(msg)
	io.Copy(packet.Buffer, reader)
	g.Game.SendToAll(packet)
}

func (g *GameAnalyser) SendToAllLen(msgTp uint8, reader *bytes.Reader, length int64) {
	if length == 0 {
		g.Game.SendToAll(NewGamePacketFactory().Create(msgTp))
		return
	}
	_, _ = reader.Seek(length, io.SeekCurrent)

	g.SendToAll(msgTp, reader)
}

func (g *GameAnalyser) SendToPlayer(msgTp uint8, reader *bytes.Reader, player int) {
	if player != 0 && player != 1 {
		return
	}
	packet := NewGamePacketFactory().Create(msgTp)
	_, _ = io.Copy(packet.Buffer, reader)
	g.Game.CurPlayers[player].Send(packet)
}

func (g *GameAnalyser) SendToTeam(msgTp uint8, reader *bytes.Reader, player int) {
	if player != 0 && player != 1 {
		return
	}

	packet := NewGamePacketFactory().Create(msgTp)
	_, _ = io.Copy(packet.Buffer, reader)
	g.Game.SendToTeam(packet, player)
}

func (g *GameAnalyser) SendToOpponentTeam(packet io.Reader, player int) {
	if player != 0 && player != 1 {
		return
	}
	g.Game.SendToTeam(packet, 1-player)

}
