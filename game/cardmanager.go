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
	tx := DB.Exec("SELECT id, ot, alias, setcode, type, level, race, attribute, atk, def FROM datas")
	if tx.Error != nil {
		return nil, tx.Error
	}
	var cardList []*Card
	err = tx.Scan(&cardList).Error
	if err != nil {
		return nil, err
	}
	var manager = &CardManager{
		cards: make(map[int32]*Card, len(cardList)),
	}
	for i := range cardList {
		manager.cards[cardList[i].Id] = cardList[i]
	}
	GlobalCardManager = manager
	return manager, nil
}
func (m *CardManager) GetCard(id int32) *Card {
	return m.cards[id]
}
