package game

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type CardManager struct {
	cards map[int32]*Card
}

var (
	GlobalCardManager *CardManager
	DB                *gorm.DB
)

func InitCardManager(file string) (*CardManager, error) {
	var err error
	DB, err = gorm.Open(sqlite.Open(file))
	if err != nil {
		return nil, err
	}
	var cardList []*Datas
	tx := DB.Select("id", "ot", "alias", "setcode", "type", "level", "race", "attribute", "atk", "def").Find(&cardList)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var manager = &CardManager{
		cards: make(map[int32]*Card, len(cardList)),
	}
	for i := range cardList {
		card := cardList[i].Card()
		manager.cards[card.Id] = card
	}
	GlobalCardManager = manager
	return manager, nil
}
func (m *CardManager) GetCard(id int32) *Card {
	return m.cards[id]
}
