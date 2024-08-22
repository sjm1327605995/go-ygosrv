package game

type Banlist struct {
	BannedIds      []int
	LimitedIds     []int
	SemiLimitedIds []int
	Hash           uint32
}

func NewBanlist() *Banlist {
	return &Banlist{
		BannedIds:      []int{},
		LimitedIds:     []int{},
		SemiLimitedIds: []int{},
		Hash:           0x7dfcee6a,
	}
}

func (b *Banlist) GetQuantity(cardId int) int {
	for _, id := range b.BannedIds {
		if id == cardId {
			return 0
		}
	}
	for _, id := range b.LimitedIds {
		if id == cardId {
			return 1
		}
	}
	for _, id := range b.SemiLimitedIds {
		if id == cardId {
			return 2
		}
	}
	return 3
}

func (b *Banlist) Add(cardId int, quantity int) {
	if quantity < 0 || quantity > 2 {
		return
	}
	switch quantity {
	case 0:
		b.BannedIds = append(b.BannedIds, cardId)
	case 1:
		b.LimitedIds = append(b.LimitedIds, cardId)
	case 2:
		b.SemiLimitedIds = append(b.SemiLimitedIds, cardId)
	}
	code := uint32(cardId)
	b.Hash = b.Hash ^ ((code << 18) | (code >> 14)) ^ ((code << (27 + uint32(quantity))) | (code >> (5 - uint32(quantity))))
}
