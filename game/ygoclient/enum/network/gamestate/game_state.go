package gamestate

type GameState uint8

const (
	Lobby GameState = iota
	Hand
	Starting
	Duel
	End
	Side
)
