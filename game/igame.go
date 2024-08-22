package game

import (
	"bytes"
	"encoding/binary"
	"fmt"
	set "github.com/duke-git/lancet/v2/datastructure/set"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/card"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/msg"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/network/gamestate"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/network/stoc"
	player2 "github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/player"
	"github.com/sjm1327605995/go-ygosrv/utils"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"math/rand"
	"os"
	"time"
)

type IGame interface {
	SetRules(packet *utils.BinaryReader) error
	AddPlayer(player *Player) error
	MoveToDuelist(player *Player) error
	MoveToObserver(player *Player) error
	RemovePlayer(player *Player) error
	SetReady(player *Player, where bool) error
	StartDuel(player *Player) error
	Surrender(player *Player, reason int, force bool) error
	Chat(player *Player, msg string) error
	HandResult(player *Player, result uint8) error
	TpResult(player *Player, result bool) error
	CustomMessage(player *Player, msg string) error
	KickPlayer(player *Player, pos uint8) error
	RefreshAll()
	RefreshAllObserver(observer *Player) error
	RefreshMonsters(player int, observer *Player) error
	RefreshSpells(player int, observer *Player) error
	RefreshHand(player int, observer *Player) error
	RefreshGrave(player int, observer *Player) error
	RefreshRemoved(player int, observer *Player) error
	RefreshExtra(player int, observer *Player) error
	WritePublicCards(update *utils.BinaryReader, result []byte) error
	RefreshSingle(player int, location int, sequence int) error
	WaitForResponse() int
}

const (
	DEFAULT_LIFEPOINTS = 8000
	DEFAULT_START_HAND = 5
	DEFAULT_DRAW_COUNT = 1
	DEFAULT_TIMER      = 240
)

type IBaseGame struct {
	Banlist        *Banlist
	Mode           int8
	Region         int8
	MasterRule     int8
	StartLp        int32
	StartHand      int32
	DrawCount      int32
	Timer          int16
	EnablePriority bool
	NoCheckDeck    bool
	NoShuffleDeck  bool
	IsMatch        bool
	IsTag          bool
	IsTpSelect     bool

	State         gamestate.GameState
	SideTimer     time.Time
	TpTimer       time.Time
	RpsTimer      time.Time
	TurnCount     int
	CurrentPlayer int
	LifePoints    [2]int32

	Players    []*Player
	CurPlayers []*Player
	IsReady    []bool
	Observers  set.Set[*Player]
	HostPlayer *Player

	Replay       *Replay
	Winner       int
	MatchResults [3]int
	MatchReasons [3]int
	DuelCount    int
	//CoreServer _server
	_duel         *Duel
	_analyser     *GameAnalyser
	_handResult   [2]int
	_startplayer  int
	_lastresponse int
	_timelimit    [2]int16
	_time         *time.Time
	_matchKill    bool

	//public event Action<object, EventArgs> OnNetworkReady;
	//public event Action<object, EventArgs> OnNetworkEnd;
	//public event Action<object, EventArgs> OnGameStart;
	//public event Action<object, EventArgs> OnGameEnd;
	//public event Action<object, EventArgs> OnDuelEnd;
	//public event Action<object, PlayerEventArgs> OnPlayerJoin;
	//public event Action<object, PlayerEventArgs> OnPlayerLeave;
	//public event Action<object, PlayerMoveEventArgs> OnPlayerMove;
	//public event Action<object, PlayerEventArgs> OnPlayerReady;
	//public event Action<object, PlayerChatEventArgs> OnPlayerChat;
}

func NewGame() *IBaseGame {
	g := new(IBaseGame)
	viper.SetDefault("Rule", -1)
	viper.SetDefault("MainDeckMinSize", 40)
	viper.SetDefault("MainDeckMaxSize", 60)
	viper.SetDefault("ExtraDeckMaxSize", 15)
	viper.SetDefault("SideDeckMaxSize", 15)
	viper.SetDefault("StartLp", DEFAULT_LIFEPOINTS)
	viper.SetDefault("StartHand", DEFAULT_START_HAND)
	viper.SetDefault("DrawCount", DEFAULT_DRAW_COUNT)
	viper.SetDefault("GameTimer", DEFAULT_TIMER)
	//  State = GameState.Lobby;
	g.Mode = int8(viper.GetInt("Mode"))
	g.Region = int8(viper.GetInt("Rule"))
	if g.Region != -1 {
		slog.Error("'Rule' is deprecated, please use 'Region' instead.")
	} else {
		g.Region = int8(viper.GetInt("Region"))
	}
	g.MasterRule = int8(viper.GetInt("MasterRule"))
	g.IsMatch = g.Mode == 1
	g.IsTag = g.Mode == 2
	/*g.CurrentPlayer = 0
	g.LifePoints = make([]int, 2)
	*/

	g.Players = make([]*Player, IFElse(g.IsTag, 4, 2))
	g.CurPlayers = make([]*Player, 2)
	g.IsReady = make([]bool, IFElse(g.IsTag, 4, 2))
	//g._handResult = make([]int, 2)
	//g._timelimit= make([]int, 2)
	g.Winner = -1
	g.Observers = set.New[*Player]()
	lfList := viper.GetInt("Banlist")
	if lfList >= 0 && lfList < len(BanListManager.Banlists) {
		b := BanListManager.Banlists[lfList]
		g.Banlist = &b
	}
	g.StartLp = viper.GetInt32("StartLp")
	g.LifePoints[0], g.LifePoints[1] = g.StartLp, g.StartLp

	g.StartHand = viper.GetInt32("StartHand")
	g.DrawCount = viper.GetInt32("DrawCount")
	g.EnablePriority = viper.GetBool("EnablePriority")
	g.NoCheckDeck = viper.GetBool("NoCheckDeck")
	g.NoShuffleDeck = viper.GetBool("NoShuffleDeck")
	g.Timer = int16(viper.GetInt("GameTimer"))

	// _server = server;
	g._analyser = &GameAnalyser{Game: g}
	return g
}
func (g *IBaseGame) SetRules(packet *bytes.Reader) error {
	var (
		lfList                                        uint32
		enablePriorityInt, noCheckDeck, noShuffleDeck int8
	)
	err := Read(packet, &lfList, g.Region, g.MasterRule, g.Mode, &enablePriorityInt, &noCheckDeck, &noShuffleDeck)
	if err != nil {
		return err
	}
	g.IsMatch = g.Mode == 1
	g.IsTag = g.Mode == 2
	g.IsReady = make([]bool, IFElse(g.IsTag, 4, 2))

	g.Players = make([]*Player, IFElse(g.IsTag, 4, 2))
	g.Players = make([]*Player, IFElse(g.IsTag, 4, 2))
	g.EnablePriority = enablePriorityInt > 0
	g.NoCheckDeck = noCheckDeck > 0
	g.NoShuffleDeck = noShuffleDeck > 0
	//C++ padding: 5 bytes + 3 bytes = 8 bytes
	packet.Seek(3, io.SeekCurrent)
	var liftPoints int32
	g.LifePoints[0] = liftPoints
	g.LifePoints[1] = liftPoints
	err = Read(packet, &liftPoints, g.StartHand, g.DrawCount, g.Timer)
	if err != nil {
		return err
	}
	return nil
}
func Read(reader io.Reader, val ...interface{}) error {
	for i := range val {
		err := binary.Read(reader, binary.LittleEndian, val[i])
		if err != nil {
			return err
		}
	}
	return nil
}
func (g *IBaseGame) Start() error {
	//if (OnNetworkReady != null)
	//{
	//	OnNetworkReady(this, EventArgs.Empty);
	//}
	return nil
}
func (g *IBaseGame) Stop() error {
	//if (OnNetworkEnd != null)
	//{
	//	OnNetworkEnd(this, EventArgs.Empty);
	//}
	return nil
}
func (g *IBaseGame) SendToAll(reader io.ReadSeeker) error {
	g.SendToPlayers(reader)

	g.SendToObservers(reader)
	return nil
}
func (g *IBaseGame) SendToAllButPlayer(reader io.ReadSeeker, player *Player) error {

	for i := range g.Players {
		if g.Players[i] != nil && !g.Players[i].Equals(player) {
			g.Players[i].Send(reader)
		}
	}
	g.Observers.Iterate(func(player *Player) {
		if !player.Equals(player) {
			player.Send(reader)
		}

	})
	return nil
}
func (g *IBaseGame) SendToAllButInt(reader io.ReadSeeker, except int) error {
	if except < len(g.CurPlayers) {
		g.SendToAllButPlayer(reader, g.CurPlayers[except])
	} else {
		g.SendToAll(reader)
	}

	return nil
}

func (g *IBaseGame) SendToPlayers(packet io.ReadSeeker) {

	for _, player := range g.Players {
		if player != nil {
			player.Send(packet)

		}
	}
}

func (g *IBaseGame) SendToObservers(packet io.ReadSeeker) {
	g.Observers.Iterate(func(player *Player) {
		player.Send(packet)
	})

}

func (g *IBaseGame) SendToTeam(packet io.ReadSeeker, team int) {

	if !g.IsTag {

		g.Players[team].Send(packet)
	} else if team == 0 {
		g.Players[0].Send(packet)
		g.Players[1].Send(packet)
	} else {
		g.Players[2].Send(packet)
		g.Players[3].Send(packet)
	}
}

func (g *IBaseGame) AddPlayer(player *Player) error {
	if g.State != gamestate.Lobby {
		player.Type = player2.Observer
		if g.State != gamestate.End {
			err := g.SendJoinGame(player)
			if err != nil {
				return err
			}
			err = player.SendTypeChange()
			if err != nil {
				return err
			}
			err = player.Send(NewGamePacketFactory().Create(stoc.DuelStart))
			if err != nil {
				return err
			}
			g.Observers.Add(player)
			if g.State == gamestate.Duel {
				g.InitNewSpectator(player)
			}
		}
		//if g.OnPlayerJoin != null {
		//	g.OnPlayerJoin(this, new
		//	PlayerEventArgs(player))
		//}
		return nil
	}

	if g.HostPlayer == nil {
		g.HostPlayer = player
	}

	pos := g.GetAvailablePlayerPos()
	if pos != -1 {

		enter := NewGamePacketFactory().Create(stoc.HsPlayerEnter)
		enter.WriteUnicode(player.Name, 20)
		enter.Write(uint8(pos))
		//padding
		enter.Write(uint8(0))
		g.SendToAll(enter)

		g.Players[pos] = player
		g.IsReady[pos] = false
		player.Type = pos
	} else {

		watch := NewGamePacketFactory().Create(stoc.HsWatchChange)
		watch.Write(int16(g.Observers.Size() + 1))
		g.SendToAll(watch)

		player.Type = player2.Observer

		g.Observers.Add(player)
	}

	g.SendJoinGame(player)
	player.SendTypeChange()
	//2b0020320032003200000031000000310000000000000000000000800d00c0597f0000cccccccccccc00000000
	//2c000120320032003200000000000000000000000000000000000000000000000000000000000000000000000000
	for i := 0; i < len(g.Players); i++ {
		if g.Players[i] != nil {

			enter := NewGamePacketFactory().Create(stoc.HsPlayerEnter)
			enter.WriteUnicode(g.Players[i].Name, 20)
			enter.Write(byte(i))
			//padding
			enter.Write(byte(0))
			player.Send(enter)

			if g.IsReady[i] {

				change := NewGamePacketFactory().Create(stoc.HsPlayerChange)
				change.Write((byte)((i << 4) + player2.Ready))
				player.Send(change)
			}
		}
	}

	if g.Observers.Size() > 0 {

		nwatch := NewGamePacketFactory().Create(stoc.HsWatchChange)
		nwatch.Write(int16(g.Observers.Size()))
		player.Send(nwatch)
	}

	//if g.OnPlayerJoin != null {
	//	OnPlayerJoin(this, new
	//	PlayerEventArgs(player))
	//}
	return nil
}

func (g *IBaseGame) RemovePlayer(player *Player) error {
	if player.Equals(g.HostPlayer) && g.State == gamestate.Lobby {
		//_server.Stop()
		return nil
	}

	if player.Type == player2.Observer {
		g.Observers.Delete(player)
		if g.State == gamestate.Lobby {
			nwatch := NewGamePacketFactory().Create(stoc.HsWatchChange)
			nwatch.Write(int16(g.Observers.Size()))
			g.SendToAll(nwatch)
		}
		player.Disconnect()
	} else if g.State == gamestate.Lobby {
		g.Players[player.Type] = nil
		g.IsReady[player.Type] = false

		change := NewGamePacketFactory().Create(stoc.HsPlayerChange)
		change.Write((byte)(player.Type<<4) + (player2.Leave))
		g.SendToAll(change)
		player.Disconnect()
	} else {
		g.Surrender(player, 4, true)
	}

	//if g.OnPlayerLeave != null {
	//	g.OnPlayerLeave(this, new
	//	PlayerEventArgs(player))
	//}
	return nil
}

func (g *IBaseGame) MoveToDuelist(player *Player) error {
	if g.State != gamestate.Lobby {
		return nil
	}
	pos := g.GetAvailablePlayerPos()
	if pos == -1 {
		return nil
	}
	//oldType := player.Type

	if player.Type != player2.Observer {
		if !g.IsTag || g.IsReady[player.Type] {
			return nil
		}

		pos = (player.Type + 1) % 4
		for g.Players[pos] != nil {
			pos = (pos + 1) % 4

			change := NewGamePacketFactory().Create(stoc.HsPlayerChange)
			change.Write((byte)((player.Type << 4) + pos))
			g.SendToAll(change)

			g.Players[player.Type] = nil
			g.Players[pos] = player
			player.Type = pos
			player.SendTypeChange()
		}

	} else {
		g.Observers.Delete(player)
		g.Players[pos] = player
		player.Type = pos

		enter := NewGamePacketFactory().Create(stoc.HsPlayerEnter)
		enter.WriteUnicode(player.Name, 20)
		enter.WriteByte(byte(pos))
		//padding
		enter.WriteByte(0)
		g.SendToAll(enter)

		nwatch := NewGamePacketFactory().Create(stoc.HsWatchChange)
		nwatch.Write(uint16(g.Observers.Size()))
		g.SendToAll(nwatch)

		player.SendTypeChange()
	}
	//if g.OnPlayerMove != nil {
	//	g.OnPlayerMove(this, new
	//	PlayerMoveEventArgs(player, oldType))
	//}
	return nil
}

func (g *IBaseGame) MoveToObserver(player *Player) error {
	if g.State != gamestate.Lobby {
		return nil
	}
	if player.Type == player2.Observer {
		return nil
	}
	if g.IsReady[player.Type] {
		return nil
	}

	//oldType := player.Type

	g.Players[player.Type] = nil
	g.IsReady[player.Type] = false

	g.Observers.Add(player)

	change := NewGamePacketFactory().Create(stoc.HsPlayerChange)
	change.Write((byte)((player.Type << 4) + (player2.Observe)))
	g.SendToAll(change)

	player.Type = player2.Observer
	player.SendTypeChange()

	//if g.OnPlayerMove != nil {
	//    g.OnPlayerMove(this, new
	//    PlayerMoveEventArgs(player, oldType))
	//}
	return nil
}

func (g *IBaseGame) Chat(player *Player, msg string) error {
	packet := NewGamePacketFactory().Create(stoc.Chat)
	packet.Write(int16(player.Type))
	if player.Type == player2.Observer {
		fullMsg := fmt.Sprintf("[%s]: %s", player.Name, msg)
		g.CustomMessage(player, fullMsg)
	} else {
		packet.WriteUnicode(msg, len(msg)+1)
		g.SendToAllButPlayer(packet, player)
	}
	//	        if (OnPlayerChat != null)
	//	        {
	//	            OnPlayerChat(this, new PlayerChatEventArgs(player, msg));
	//	        }

	return nil
}

func (g *IBaseGame) CustomMessage(player *Player, msg string) error {
	packet := NewGamePacketFactory().Create(stoc.Chat)
	packet.Write(int16(player2.Yellow))

	packet.WriteUnicode(msg, len(msg)+1)
	return g.SendToAllButPlayer(packet, player)
}

func (g *IBaseGame) SetReady(player *Player, ready bool) error {
	if g.State != gamestate.Lobby {
		return nil
	}
	if player.Type == player2.Observer {
		return nil
	}
	if g.IsReady[player.Type] == ready {
		return nil
	}
	if ready {
		ocg := g.Region == 0 || g.Region == 2
		tcg := g.Region == 1 || g.Region == 2
		result := int32(1)
		if player.Deck != nil {
			result = IFElse(g.NoCheckDeck, 0, player.Deck.Check(g.Banlist, ocg, tcg))
		}
		if result != 0 {
			rechange := NewGamePacketFactory().Create(stoc.HsPlayerChange)
			rechange.Write((byte)((player.Type << 4) + (player2.NotReady)))
			player.Send(rechange)
			errMsg := NewGamePacketFactory().Create(stoc.ErrorMsg)
			errMsg.Write(byte(2)) // ErrorMsg.DeckErr
			//	                // C++ padding: 1 byte + 3 bytes = 4 bytes
			for i := 0; i < 3; i++ {
				errMsg.Write(byte(0))
			}
			errMsg.Write(result)
			player.Send(rechange)
		}
		g.IsReady[player.Type] = ready
	}

	change := NewGamePacketFactory().Create(stoc.HsPlayerChange)
	change.Write(byte((player.Type << 4) + IFElse(ready, player2.Ready, player2.NotReady)))
	return g.SendToAll(change)

	//	        if (OnPlayerReady != null)
	//	        {
	//	            OnPlayerReady(this, new PlayerEventArgs(player));
	//	        }
}

func (g *IBaseGame) KickPlayer(player *Player, pos uint8) error {
	if g.State != gamestate.Lobby {
		return nil
	}

	if int(pos) >= len(g.Players) || !player.Equals(g.HostPlayer) || player.Equals(g.Players[pos]) || g.Players[pos] == nil {
		return nil
	}

	return g.RemovePlayer(g.Players[pos])
}

func (g *IBaseGame) StartDuel(player *Player) error {

	if g.State != gamestate.Lobby {
		return nil
	}

	if !player.Equals(g.HostPlayer) {
		return nil
	}

	for i := 0; i < len(g.Players); i++ {
		if !g.IsReady[i] {
			return nil
		}

		if g.Players[i] == nil {
			return nil
		}

	}

	g.State = gamestate.Hand
	g.SendToAll(NewGamePacketFactory().Create(stoc.DuelStart))

	g.SendHand()

	//if (OnGameStart != null)
	//{
	//    OnGameStart(this, EventArgs.Empty);
	//}
	return nil
}

func (g *IBaseGame) HandResult(player *Player, result uint8) error {
	if g.State != gamestate.Hand {
		return nil
	}
	if player.Type == player2.Observer {
		return nil
	}
	if result < 1 || result > 3 {
		return nil
	}
	if g.IsTag && player.Type != 0 && player.Type != 2 {
		return nil
	}
	tp := player.Type
	if g.IsTag && player.Type == 2 {
		tp = 1
	}

	if g._handResult[tp] != 0 {
		return nil
	}
	g._handResult[tp] = int(result)
	if g._handResult[0] != 0 && g._handResult[1] != 0 {
		packet := NewGamePacketFactory().Create(stoc.HandResult)
		packet.Write(byte(g._handResult[0]))
		packet.Write(byte(g._handResult[1]))
		g.SendToTeam(packet, 0)
		g.SendToObservers(packet)

		packet = NewGamePacketFactory().Create(stoc.HandResult)
		packet.Write(byte(g._handResult[0]))
		packet.Write(byte(g._handResult[1]))
		g.SendToTeam(packet, 1)

		if g._handResult[0] == g._handResult[1] {
			g._handResult[0] = 0
			g._handResult[1] = 0
			g.SendHand()
			return nil
		}
		if (g._handResult[0] == 1 && g._handResult[1] == 2) ||
			(g._handResult[0] == 2 && g._handResult[1] == 3) ||
			(g._handResult[0] == 3 && g._handResult[1] == 1) {
			g._startplayer = IFElse(g.IsTag, 2, 1)
		} else {
			g._startplayer = 0
		}

		g.State = gamestate.Starting
		g.Players[g._startplayer].Send(NewGamePacketFactory().Create(stoc.SelectTp))
		g.TpTimer = time.Now().UTC()
	}
	return nil
}

func (g *IBaseGame) TpResult(player *Player, result bool) error {
	if g.State != gamestate.Starting {
		return nil
	}
	if player.Type != g._startplayer {
		return nil
	}
	opt := int(g.MasterRule) << 16
	if g.EnablePriority {
		opt += 0x08
	}

	if g.NoShuffleDeck {
		opt += 0x10
	}

	if g.IsTag {
		opt += 0x20
	}

	if result && player.Type == IFElse(g.IsTag, 2, 1) || !result && player.Type == 0 {
		opt += 0x80
	}

	g.CurPlayers[0] = g.Players[0]
	g.CurPlayers[1] = g.Players[IFElse(g.IsTag, 2, 1)]
	g.State = gamestate.Duel
	seed := cast.ToUint32(os.Getenv("TickCount"))

	g._duel = NewDuel(seed)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	g._duel.SetAnalyzer(g._analyser.Analyser)
	g._duel.SetErrorHandler(g.HandleError)

	g._duel.InitPlayers(g.StartLp, g.StartHand, g.DrawCount)

	//g.Replay = NewReplay(seed, g.IsTag)
	//g.Replay.WriteUnicode(g.Players[0].Name, 20)
	//g.Replay.WriteUnicode(g.Players[1].Name, 20)
	//if g.IsTag {
	//	g.Replay.WriteUnicode(g.Players[2].Name, 20)
	//	g.Replay.WriteUnicode(g.Players[3].Name, 20)
	//}
	//g.Replay.Writer.Write(g.StartLp)
	//g.Replay.Writer.Write(g.StartHand)
	//g.Replay.Writer.Write(g.DrawCount)
	//g.Replay.Writer.Write(opt)

	for i := 0; i < len(g.Players); i++ {
		dplayer := g.Players[i]
		pid := i
		if g.IsTag {
			pid = IFElse(i >= 2, 1, 0)
		}
		if !g.NoShuffleDeck {
			cards := g.ShuffleCards(r, dplayer.Deck.Main.Data())
			//g.Replay.Writer.Write(len(cards))
			for _, id := range cards {
				if g.IsTag && (i == 1 || i == 3) {
					g._duel.AddTagCard(id, pid, card.Deck)
				} else {
					g._duel.AddCard(id, pid, card.Deck)
				}
				//g.Replay.Writer.Write(id)
			}
		} else {
			//g.Replay.Writer.Write(len(dplayer.Deck.Main))
			cards := dplayer.Deck.Main.Data()

			for j := len(cards) - 1; j >= 0; j-- {
				id := cards[j]
				if g.IsTag && (i == 1 || i == 3) {
					g._duel.AddTagCard(id, pid, card.Deck)
				} else {
					g._duel.AddCard(id, pid, card.Deck)
				}
				//	g.Replay.Writer.Write(id)
			}
		}
		//g.Replay.Writer.Write(len(dplayer.Deck.Extra))
		dplayer.Deck.Extra.ForEach(func(id int32) {
			if g.IsTag && (i == 1 || i == 3) {
				g._duel.AddTagCard(id, pid, card.Extra)
			} else {
				g._duel.AddCard(id, pid, card.Extra)
			}
			//g.Replay.Writer.Write(id)
		})

	}

	packet := NewGamePacketFactory().Create(msg.Start)
	packet.Write(byte(0))
	packet.Write(byte(g.MasterRule))
	packet.Write(g.StartLp)
	packet.Write(g.StartLp)
	packet.Write(int16(g._duel.QueryFieldCount(0, card.Deck)))
	packet.Write(int16(g._duel.QueryFieldCount(0, card.Extra)))
	packet.Write(int16(g._duel.QueryFieldCount(1, card.Deck)))
	packet.Write(int16(g._duel.QueryFieldCount(1, card.Extra)))

	g.SendToTeam(packet, 0)
	packet.Seek(2, io.SeekStart)
	//packet.BaseStream.Position = 2
	packet.WriteByte(1)
	g.SendToTeam(packet, 1)

	packet.Seek(2, io.SeekStart)
	packet.Write(byte(0x10))
	g.SendToObservers(packet)

	g.RefreshExtra(0, nil)
	g.RefreshExtra(1, nil)

	g._duel.Start(opt)

	g.TurnCount = 0
	g.LifePoints[0] = g.StartLp
	g.LifePoints[1] = g.StartLp
	g.TimeReset()

	g.Process()
	return nil
}

func (g *IBaseGame) Surrender(player *Player, reason int, force bool) error {
	if g.State == gamestate.End {
		return nil
	}
	if !force && g.State != gamestate.Duel {
		return nil
	}
	if player.Type == int(player2.Observer) {
		return nil
	}
	win := NewGamePacketFactory().Create(msg.Win)
	team := player.Type
	if g.IsTag {
		team = IFElse(player.Type >= 2, 1, 0)
	} else if g.State == gamestate.Hand {
		team = 1 - team
	}
	win.Write(byte(1 - team))
	win.Write(byte(reason))
	g.SendToAll(win)

	g.MatchSaveResult(1-team, reason)

	g.EndDuel(reason == 4)
	return nil
}

func (g *IBaseGame) RefreshAll() {
	g.RefreshMonsters(0, nil)
	g.RefreshMonsters(1, nil)
	g.RefreshSpells(0, nil)
	g.RefreshSpells(1, nil)
	g.RefreshHand(0, nil)
	g.RefreshHand(1, nil)
}

func (g *IBaseGame) RefreshAllObserver(observer *Player) error {
	g.RefreshMonsters(0, observer)
	g.RefreshMonsters(1, observer)
	g.RefreshSpells(0, observer)
	g.RefreshSpells(1, observer)
	g.RefreshHand(0, observer)
	g.RefreshHand(1, observer)
	g.RefreshGrave(0, observer)
	g.RefreshGrave(1, observer)
	g.RefreshExtra(0, observer)
	g.RefreshExtra(1, observer)
	g.RefreshRemoved(0, observer)
	g.RefreshRemoved(1, observer)
	return nil
}

func (g *IBaseGame) RefreshMonsters(player int, observer *Player) error {
	result := g._duel.QueryFieldCard(player, card.MonsterZone, 0xFFFFFF & ^enum.ReasonCard, false)
	g.SendToCorrectDestination(player, card.MonsterZone, result, observer)
	return nil
}

func (g *IBaseGame) RefreshSpells(player int, observer *Player) error {
	result := g._duel.QueryFieldCard(player, card.SpellZone, 0xFFFFFF & ^enum.ReasonCard, false)
	g.SendToCorrectDestination(player, card.SpellZone, result, observer)
	return nil
}

func (g *IBaseGame) RefreshHand(player int, observer *Player) error {
	result := g._duel.QueryFieldCard(player, card.Hand, 0xFFFFFF & ^enum.ReasonCard, false)
	g.SendToCorrectDestination(player, card.Hand, result, observer)
	return nil

}

func (g *IBaseGame) RefreshGrave(player int, observer *Player) error {
	result := g._duel.QueryFieldCard(player, card.Grave, 0xFFFFFF & ^enum.ReasonCard, false)
	g.SendToCorrectDestination(player, card.Grave, result, observer)
	return nil
}

func (g *IBaseGame) RefreshRemoved(player int, observer *Player) error {
	result := g._duel.QueryFieldCard(player, card.Removed, 0xFFFFFF & ^enum.ReasonCard, false)
	g.SendToCorrectDestination(player, card.Removed, result, observer)
	return nil
}

func (g *IBaseGame) RefreshExtra(player int, observer *Player) error {
	result := g._duel.QueryFieldCard(player, card.Extra, 0xFFFFFF & ^enum.ReasonCard, false)
	g.SendToCorrectDestination(player, card.Extra, result, observer)
	return nil
}

func (g *IBaseGame) SendToCorrectDestination(player int, location uint8, result []byte, observer *Player) {
	var update *GamePacketFactory

	if observer == nil {
		update = NewGamePacketFactory().Create(msg.UpdateData)
		update.Write(byte(player))
		update.Write(location)
		update.Write(result)
		g.SendToTeam(update, player)
	}

	update = NewGamePacketFactory().Create(msg.UpdateData)
	update.Write(byte(player))
	update.Write(location)
	g.WritePublicCards(update.MemoryStream, result)

	if observer == nil {
		g.SendToTeam(update, 1-player)
		g.SendToObservers(update)
	} else {
		observer.Send(update)
	}

}

func (g *IBaseGame) WritePublicCards(update io.Writer, result []byte) error {
	reader := bytes.NewBuffer(result)

	for {
		var length int32
		err := binary.Read(reader, binary.LittleEndian, &length)
		if err != nil {
			return err
		}

		if length == 4 {
			update.Write([]byte{4})
			continue
		}
		raw := reader.Next(int(length - 4))

		isFaceup := raw[11]&card.FaceUp != 0
		if isFaceup {
			update.Write([]byte{byte(length)})
			update.Write(raw)
		} else {
			update.Write([]byte{8, 0})

		}
	}
	return nil
}

func (g *IBaseGame) RefreshSingle(player int, location int, sequence int) error {
	result := g._duel.QueryCard(player, location, sequence, 0xFFFFFF & ^enum.ReasonCard, false)

	if location == card.Removed && result[15]&card.FaceDown != 0 {
		return nil
	}

	update := NewGamePacketFactory().Create(msg.UpdateCard)
	update.Write(byte(player))
	update.Write(byte(location))
	update.Write(byte(sequence))
	update.Write(result)
	g.CurPlayers[player].Send(update)

	if g.IsTag {
		if (location & card.OnField) != 0 {
			g.SendToTeam(update, player)
			if (result[15] & card.FaceUp) != 0 {
				g.SendToTeam(update, 1-player)
			}
		} else {
			g.CurPlayers[player].Send(update)
			if (location & 0x90) != 0 {
				g.SendToAllButInt(update, player)
			}
		}
	} else {
		if (location&0x90) != 0 || ((location&0x2c) != 0 && result[15]&card.FaceUp != 0) {
			g.SendToAllButInt(update, player)
		}
	}
	return nil
}

func (g *IBaseGame) WaitForResponse() int {
	g.WaitForResponseN(g._lastresponse)
	return g._lastresponse
}

func (g *IBaseGame) WaitForResponseN(player int) {
	g._lastresponse = player
	g.CurPlayers[player].State = player2.Response
	g.SendToAllButInt(NewGamePacketFactory().Create(msg.Waiting), player)
	g.TimeStart()
	packet := NewGamePacketFactory().Create(stoc.TimeLimit)
	packet.Write(byte(player))
	packet.Write(byte(0)) // Go padding
	packet.Write(int16(g._timelimit[player]))
	g.SendToPlayers(packet)
}

func (g *IBaseGame) SetResponse(resp int) error {
	//if !g.Replay.Disabled {
	//	g.Replay.Writer.Write(byte(4))
	//	g.Replay.Writer.Write(BitConverter.GetBytes(resp))
	//	g.Replay.Check()
	//}

	g.TimeStop()
	g._duel.SetResponseInt(resp)
	return nil
}

func (g *IBaseGame) SetResponseBytes(resp []byte) error {
	//if !g.Replay.Disabled {
	//	g.Replay.Writer.Write(byte(len(resp)))
	//	g.Replay.Writer.Write(resp)
	//	g.Replay.Check()
	//}

	g.TimeStop()
	g._duel.SetResponseByte(resp)
	g.Process()
	return nil
}

func (g *IBaseGame) EndDuel(force bool) {
	if g.State == gamestate.End {
		return
	}
	if g.State == gamestate.Duel {
		//	            if (!Replay.Disabled)
		//	            {
		//	                Replay.End();
		//	                byte[] replayData = Replay.GetContent();
		//	                BinaryWriter packet = GamePacketFactory.Create(stoc.Replay);
		//	                packet.Write(replayData);
		//	                SendToAll(packet);
		//	            }
		//
		//	            _duel.End();
	}
	if g.IsMatch && !force && g.MatchIsEnd() {
		g.IsReady[0] = false
		g.IsReady[1] = false
		g.State = gamestate.Side
		//	            SideTimer = DateTime.UtcNow;
		g.SendToPlayers(NewGamePacketFactory().Create(stoc.ChangeSide))
		g.SendToObservers(NewGamePacketFactory().Create(stoc.WaitingSide))
	} else {
		g.CalculateWinner()
		g.End()
	}
}
func (g *IBaseGame) End() {
	g.State = gamestate.End

	g.SendToAll(NewGamePacketFactory().Create(stoc.DuelEnd))
	//  _server.StopDelayed();

	//if g.OnGameEnd {
	//	OnGameEnd(this, EventArgs.Empty)
	//}
}
func (g *IBaseGame) TimeReset() {
	g._timelimit[0] = g.Timer
	g._timelimit[1] = g.Timer
}

func (g *IBaseGame) TimeStart() {
	t := time.Now().UTC()
	g._time = &t
}

func (g *IBaseGame) TimeStop() {
	if g._time != nil {
		elapsed := time.Now().UTC().Sub(*g._time)
		g._timelimit[g._lastresponse] -= int16(elapsed.Seconds())
		if g._timelimit[g._lastresponse] < 0 {
			g._timelimit[g._lastresponse] = 0
		}
		g._time = nil
	}
}

func (g *IBaseGame) TimeTick() {
	if g.State == gamestate.Duel {
		if g._time != nil {
			elapsed := time.Now().UTC().Sub(*g._time)
			if int16(elapsed.Seconds()) > g._timelimit[g._lastresponse] {
				g.Surrender(g.CurPlayers[g._lastresponse], 3, false)
			}
		}
	}

	if g.State == gamestate.Side {
		elapsed := time.Now().UTC().Sub(g.SideTimer)

		if elapsed.Milliseconds() >= 120000 {
			if !g.IsReady[0] && !g.IsReady[1] {
				g.EndDuel(true)
				return
			}

			g.Surrender(IFElse(!g.IsReady[0], g.Players[0], g.Players[1]), 3, true)
		}
	}

	if g.State == gamestate.Starting {
		if g.IsTpSelect {
			elapsed := time.Now().UTC().Sub(g.TpTimer)

			if elapsed.Milliseconds() >= 30000 {
				g.Surrender(g.CurPlayers[g._startplayer], 3, true)
			}
		}
	}

	if g.State == gamestate.Hand {
		elapsed := time.Now().UTC().Sub(g.RpsTimer)

		if int(elapsed.Milliseconds()) >= 60000 {
			if g._handResult[0] != 0 {
				g.Surrender(g.Players[IFElse(g.IsTag, 2, 1)], 3, true)
			} else if g._handResult[1] != 0 {
				g.Surrender(g.Players[0], 3, true)
			} else {
				g.EndDuel(true)
			}
		}
	}
}

func (d *IBaseGame) MatchSaveResult(player, reason int) {
	if player < 2 {
		d._startplayer = 1 - player
	} else {
		d._startplayer = 1 - d._startplayer
	}
	d.MatchResults[d.DuelCount] = player
	d.MatchReasons[d.DuelCount] = reason
	d.DuelCount++

	//if d.OnDuelEnd != nil {
	//	d.OnDuelEnd(d, EventArgsEmpty)
	//}
}

func (d *IBaseGame) MatchKill() {
	d._matchKill = true
}

func (d *IBaseGame) MatchIsEnd() bool {
	if d._matchKill {
		return true
	}
	wins := [3]int{}
	for i := 0; i < d.DuelCount; i++ {
		wins[d.MatchResults[i]]++
	}
	return wins[0] == 2 || wins[1] == 2 || wins[0]+wins[1]+wins[2] == 3
}

func (d *IBaseGame) MatchSide() {
	if d.IsReady[0] && d.IsReady[1] {
		d.State = gamestate.Starting
		d.IsTpSelect = true
		d.TpTimer = time.Now().UTC()
		d.TimeReset()
		d.Players[d._startplayer].Send(NewGamePacketFactory().Create(stoc.SelectTp))
	}
}

func (d *IBaseGame) GetAvailablePlayerPos() int {
	for i := 0; i < len(d.Players); i++ {
		if d.Players[i] == nil {
			return i
		}
	}
	return -1
}

func (d *IBaseGame) SendHand() {
	d.RpsTimer = time.Now().UTC()
	hand := NewGamePacketFactory().Create(stoc.SelectHand)
	if d.IsTag {
		d.Players[0].Send(hand)
		d.Players[2].Send(hand)
	} else {
		d.SendToPlayers(hand)
	}
}

func (d *IBaseGame) Process() {
	result := d._duel.Process()
	switch result {
	case -1:
		d.EndDuel(true)
	case 2: // Game finished
		d.EndDuel(false)
	}
}

func (g *IBaseGame) SendJoinGame(player *Player) error {
	join := NewGamePacketFactory().Create(stoc.JoinGame)
	if g.Banlist == nil {
		join.Write(uint32(0))
	} else {
		join.Write(g.Banlist.Hash)
	}
	join.Write(byte(g.Region))
	join.Write(byte(g.Mode))
	join.Write(byte(g.MasterRule))
	join.Write(g.NoCheckDeck)
	join.Write(g.NoShuffleDeck)
	// Go padding: 5 bytes + 3 bytes = 8 bytes
	for i := 0; i < 3; i++ {
		join.Write(byte(0))
	}
	join.Write(g.StartLp)
	join.Write(byte(g.StartHand))
	join.Write(byte(g.DrawCount))
	join.Write(int16(g.Timer))
	err := player.Send(join)
	if err != nil {
		return err
	}
	if g.State != gamestate.Lobby {
		g.SendDuelingPlayers(player)
	}
	return nil
}

func (g *IBaseGame) SendDuelingPlayers(player *Player) error {
	for i := 0; i < len(g.Players); i++ {
		enter := NewGamePacketFactory().Create(stoc.HsPlayerEnter)
		enter.WriteUnicode(g.Players[i].Name, 20)
		enter.Write(byte(i))
		//padding
		enter.Write(byte(0))
		err := player.Send(enter)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (g *IBaseGame) InitNewSpectator(player *Player) {
	packet := NewGamePacketFactory().Create(msg.Start)
	packet.Write(byte(0x10)).
		Write(byte(g.MasterRule)).
		Write(g.LifePoints[0]).
		Write(g.LifePoints[1]).
		Write(int16(0)) // deck .Write(int16(0)) // extra .Write(int16(0)) // deck .Write(int16(0))  // extra
	player.Send(packet)

	turn := NewGamePacketFactory().Create(msg.NewTurn)
	turn.Write(byte(0))
	player.Send(turn)
	if g.CurrentPlayer == 1 {
		turn = NewGamePacketFactory().Create(msg.NewTurn)
		turn.Write(byte(0))
		player.Send(turn)
	}

	reload := NewGamePacketFactory().Create(msg.ReloadField)
	fieldInfo := g._duel.QueryFieldInfo()
	reload.WriteLen(fieldInfo[1:], len(fieldInfo)-1)
	player.Send(reload)

	g.RefreshAllObserver(player)
}

func (g *IBaseGame) HandleError(error string) {
	packet := make([]byte, 2+len(error)+1)
	binary.BigEndian.PutUint16(packet, uint16(0)) // Assuming PlayerType.Observer is 0
	for i, r := range error {
		packet[2+i] = byte(r)
	}
	packet[len(error)+2] = 0 // Null-terminate the string

	// Assuming SendToAll is a function that sends the packet to all clients
	g.SendToAll(bytes.NewReader(packet))

	filename := fmt.Sprintf("lua_%s.txt", time.Now().Format("2006-01-02_15-04-05"))
	err := os.WriteFile(filename, []byte(error), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func (g *IBaseGame) ShuffleCards(rand *rand.Rand, cards []int32) []int32 {
	shuffled := make([]int32, len(cards))
	copy(shuffled, cards)
	for i := len(shuffled) - 1; i > 0; i-- {
		pos := rand.Intn(i + 1)
		shuffled[i], shuffled[pos] = shuffled[pos], shuffled[i]
	}
	return shuffled
}

func (g *IBaseGame) CalculateWinner() {
	winner := -1
	if g.DuelCount > 0 {
		if !g._matchKill && g.DuelCount != 1 {
			wins := make([]int, 3)
			for i := 0; i < g.DuelCount; i++ {
				wins[g.MatchResults[i]]++
			}
			if wins[0] > wins[1] {
				winner = 0
			} else if wins[1] > wins[0] {
				winner = 1
			} else {
				winner = 2
			}
		} else {
			winner = g.MatchResults[g.DuelCount-1]
		}
	}
	g.Winner = winner
}

func IFElse[T any](where bool, val1, val2 T) T {
	if where {
		return val1
	}
	return val2
}
