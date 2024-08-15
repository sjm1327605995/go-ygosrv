package game

import (
	"encoding/base64"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/network/gamestate"
	"log"
	"math/rand"
	"strings"
)

type GameManager struct {
	rooms map[string]*GameRoom
}

func NewGameManager() *GameManager {
	return &GameManager{
		rooms: make(map[string]*GameRoom),
	}
}

func (gm *GameManager) CreateOrGetGame(config *GameConfig) *GameRoom {
	if room, exists := gm.rooms[config.Name]; exists {
		return room
	}
	return gm.CreateRoom(config)
}

func (gm *GameManager) GetGame(name string) *GameRoom {
	if room, exists := gm.rooms[name]; exists {
		return room
	}
	return nil
}

//func (gm *GameManager) GetRandomGame(filter int) *GameRoom {
//	filteredRooms := []*GameRoom{}
//	for _, room := range gm.rooms {
//		if room.Game.State == gamestate.Lobby && (filter == -1 || room.Game.Config.Rule == filter) {
//			filteredRooms = append(filteredRooms, room)
//		}
//	}
//
//	if len(filteredRooms) == 0 {
//		return nil
//	}
//
//	return filteredRooms[rand.Intn(len(filteredRooms))]
//}

func (gm *GameManager) SpectateRandomGame() *GameRoom {
	filteredRooms := []*GameRoom{}
	for _, room := range gm.rooms {
		if room.Game.State != gamestate.Lobby {
			filteredRooms = append(filteredRooms, room)
		}
	}

	if len(filteredRooms) == 0 {
		return nil
	}

	return filteredRooms[rand.Intn(len(filteredRooms))]
}

func (gm *GameManager) CreateRoom(config *GameConfig) *GameRoom {
	room := NewGameRoom(config)
	gm.rooms[config.Name] = room
	log.Println("Game++")
	return room
}

func (gm *GameManager) HandleRooms() {
	toRemove := []string{}
	for key, room := range gm.rooms {
		if room.IsOpen {
			room.HandleGame()
		} else {
			toRemove = append(toRemove, key)
		}
	}

	for _, room := range toRemove {
		delete(gm.rooms, room)
		log.Println("Game--")
	}
}

func (gm *GameManager) GameExists(name string) bool {
	_, exists := gm.rooms[name]
	return exists
}

func (gm *GameManager) RandomRoomName() string {
	for {
		guidBytes := make([]byte, 16)
		rand.Read(guidBytes)
		guidString := base64.StdEncoding.EncodeToString(guidBytes)
		guidString = strings.ReplaceAll(guidString, "=", "")
		guidString = strings.ReplaceAll(guidString, "+", "")
		roomName := guidString[:5]
		if !gm.GameExists(roomName) {
			return roomName
		}
	}
}
