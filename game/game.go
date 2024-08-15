package game

import "github.com/panjf2000/gnet/v2"

type GameClient struct {
	YgoClient  gnet.Conn
	Username   string
	Deck       string
	DeckFile   string
	Dialog     string
	Hand       int
	Debug      bool
	Chat       bool
	ServerHost string
	severPort  int
	roomInfo   string
	// GameBehavior _behavior
	_proVersion int16
}

func (g *GameClient) Start() {

}
