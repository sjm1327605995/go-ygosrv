package game

import (
	"github.com/sjm1327605995/go-ygocore"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/card"
)

type Card struct {
	*ygocore.CardData
	Id         int32
	Ot         int
	Alias      int32
	SetCode    int64
	Type       int
	Level      int
	LScale     int
	RScale     int
	LinkMarker int
	Attribute  int
	Race       int
	Attack     int
	Defense    int
}

func GetCard(id int32) *Card {
	return GlobalCardManager.GetCard(id)
}

// HasType CardType
func (c *Card) HasType(tp uint16) bool {
	return c.Type&int(tp) != 0
}
func (c *Card) IsExtraCard() bool {
	return c.HasType(card.Fusion) || c.HasType(card.Synchro) || c.HasType(card.Xyz) || c.HasType(card.Link)
}
