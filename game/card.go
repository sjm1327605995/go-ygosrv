package game

import (
	"github.com/sjm1327605995/go-ygocore"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/card"
)

type Card struct {
	*ygocore.CardData
	Id         int32
	Ot         uint32
	Alias      int32
	SetCode    int64
	Type       int32
	Level      uint32
	LScale     uint32
	RScale     uint32
	LinkMarker uint32
	Attribute  uint32
	Race       uint32
	Attack     int32
	Defense    int32
}
type Datas struct {
	ID        uint32 `gorm:"column:id" db:"id" json:"id" form:"id"`
	Ot        uint32 `gorm:"column:ot" db:"ot" json:"ot" form:"ot"`
	Alias     uint32 `gorm:"column:alias" db:"alias" json:"alias" form:"alias"`
	Setcode   int64  `gorm:"column:setcode" db:"setcode" json:"setcode" form:"setcode"`
	Type      int32  `gorm:"column:type" db:"type" json:"type" form:"type"`
	Atk       int32  `gorm:"column:atk" db:"atk" json:"atk" form:"atk"`
	Def       int32  `gorm:"column:def" db:"def" json:"def" form:"def"`
	Level     uint32 `gorm:"column:level" db:"level" json:"level" form:"level"`
	Race      uint32 `gorm:"column:race" db:"race" json:"race" form:"race"`
	Attribute uint32 `gorm:"column:attribute" db:"attribute" json:"attribute" form:"attribute"`
	Category  uint32 `gorm:"column:category" db:"category" json:"category" form:"category"`
}

func (d *Datas) TableName() string {
	return "datas"
}
func (d *Datas) Card() *Card {
	level := d.Level & 0xff

	var c = Card{
		Id:        int32(d.ID),
		Ot:        d.Ot,
		Alias:     int32(d.Alias),
		SetCode:   d.Setcode,
		Type:      d.Type,
		LScale:    (level >> 24) & 0xff,
		RScale:    (level >> 16) & 0xff,
		Race:      d.Race,
		Attribute: d.Attribute,
		Attack:    d.Atk,
		Defense:   d.Def,
	}
	if c.HasType(card.Link) {
		c.LinkMarker = uint32(c.Defense)
		c.Defense = 0
	}

	c.CardData = &ygocore.CardData{
		Code:       uint32(c.Id),
		Ot:         c.Ot,
		Alias:      uint32(c.Alias),
		SetCode:    uint64(c.SetCode),
		Typ:        c.Type,
		Level:      c.Level,
		Race:       c.Race,
		Attribute:  c.Attribute,
		Attack:     c.Attack,
		Defense:    c.Defense,
		Lscale:     c.LScale,
		Rscale:     c.RScale,
		LinkMarker: int32(c.LinkMarker),
	}

	return &c
}
func GetCard(id int32) *Card {
	return GlobalCardManager.GetCard(id)
}

// HasType CardType
func (c *Card) HasType(tp int32) bool {
	return c.Type&tp != 0
}
func (c *Card) IsExtraCard() bool {
	return c.HasType(card.Fusion) || c.HasType(card.Synchro) || c.HasType(card.Xyz) || c.HasType(card.Link)
}
