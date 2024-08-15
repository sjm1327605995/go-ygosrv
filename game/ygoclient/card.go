package ygoclient

import "sync"

// CardData struct holds the data for a card.
type CardData struct {
	Code      int32
	Alias     int32
	Setcode   int64
	Type      int32
	Level     int32
	Attribute int32
	Race      int32
	Attack    int32
	Defense   int32
	LScale    int32
	RScale    int32
}

// Card represents a card in the game.
type Card struct {
	Id   int
	Ot   int
	Data CardData
}

// CardsManager is a placeholder for the actual card management logic.
var CardsManager = &CardsManagerType{}

type CardsManagerType struct {
	mu sync.Mutex
}

// GetCard simulates getting a card by its ID.
func (cm *CardsManagerType) GetCard(id int) *Card {
	// Placeholder for actual card retrieval logic.
	return &Card{Id: id, Ot: 0, Data: CardData{}}
}

// Get retrieves a card by its ID.
func Get(id int) *Card {
	return CardsManager.GetCard(id)
}

// LoadCardEventHandler defines the event handler type for loading a card.
type LoadCardEventHandler func(cardId int, card *CardData)

// LoadCard is the event that gets triggered when a card is loaded.
var LoadCard LoadCardEventHandler

// OnLoadCard triggers the LoadCard event.
func OnLoadCard(cardId int, card *CardData) {
	if LoadCard != nil {
		LoadCard(cardId, card)
	}
}

// NewCard creates a new Card instance.
func NewCard(data CardData, ot int) *Card {
	return &Card{
		Id:   int(data.Code),
		Ot:   ot,
		Data: data,
	}
}
