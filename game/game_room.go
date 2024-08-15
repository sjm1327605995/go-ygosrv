package game

type GameRoom struct {
	Game         *IBaseGame
	clients      []GameClient
	IsOpen       bool
	closePending bool
}

func NewGameRoom(config *GameConfig) *GameRoom {
	return &GameRoom{
		clients: []GameClient{},
		//Game:    NewGame(this, config),
		IsOpen: true,
	}
}

func (room *GameRoom) AddClient(client GameClient) {
	room.clients = append(room.clients, client)
}

func (room *GameRoom) RemoveClient(client GameClient) {
	for i, c := range room.clients {
		if c == client {
			room.clients = append(room.clients[:i], room.clients[i+1:]...)
			break
		}
	}
}

func (room *GameRoom) Close() {
	//room.IsOpen = false
	//for _, client := range room.clients {
	//	//client.Close()
	//}
}

func (room *GameRoom) CloseDelayed() {
	//for _, client := range room.clients {
	//	client.CloseDelayed()
	//}
	room.closePending = true
}

func (room *GameRoom) HandleGame() {
	//for _, user := range room.clients {
	//	user.Tick()
	//}

	room.Game.TimeTick()

	if room.closePending && len(room.clients) == 0 {
		room.Close()
	}
}
