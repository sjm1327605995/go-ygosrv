package game

import (
	list "github.com/duke-git/lancet/v2/datastructure/list"
	"github.com/sjm1327605995/go-ygosrv/game/ygoclient/enum/card"
	"github.com/spf13/viper"
)

type Deck struct {
	Main  list.List[int32]
	Extra list.List[int32]
	Side  list.List[int32]
}

func (d *Deck) AddMain(cardId int32) {
	ACard := GetCard(cardId)
	if ACard == nil {
		return
	}
	if ACard.Type&card.Token != 0 {
		return
	}
	if ACard.IsExtraCard() {
		if d.Extra.Size() < 15 {
			d.Extra.Push(cardId)
		}
	} else {
		if d.Main.Size() < 60 {
			d.Main.Push(cardId)
		}
	}
}
func (d *Deck) AddSide(cardId int32) {
	aCard := GetCard(cardId)
	if aCard == nil {
		return
	}
	if aCard.Type&card.Token != 0 {
		return
	}
	if d.Side.Size() < 15 {
		d.Side.Push(cardId)
	}
}
func (d *Deck) Check(ban *Banlist, ocg bool, tcg bool) int32 {
	if d.Main.Size() < viper.GetInt("MainDeckMinSize") || d.Main.Size() > viper.GetInt("MainDeckMaxSize") ||
		d.Extra.Size() > viper.GetInt("ExtraDeckMaxSize") || d.Side.Size() > viper.GetInt("SideDeckMaxSize") {
		return 1
	}

	cardsMap := make(map[int32]int)
	stacks := []list.List[int32]{d.Main, d.Extra, d.Side}

	for _, stack := range stacks {
		iterator := stack.Iterator()
		for iterator.HasNext() {
			id, _ := iterator.Next()
			aCard := GetCard(id)
			d.AddToCards(cardsMap, aCard)
			if !ocg && aCard.Ot == 1 || !tcg && aCard.Ot == 2 {
				return id
			}

		}

	}

	if ban == nil {
		return 0
	}

	for k, v := range cardsMap {
		maxV := ban.GetQuantity(int(k))
		if v > maxV {
			return k
		}

	}

	return 0
}

func (d *Deck) CheckBool(deck *Deck) bool {
	if deck.Main.Size() != deck.Main.Size() || deck.Extra.Size() != d.Extra.Size() {
		return false
	}
	var (
		cards  = make(map[int32]int)
		ncards = make(map[int32]int)
	)
	stacks := []list.List[int32]{d.Main, d.Extra, d.Side}
	for _, stack := range stacks {
		iterator := stack.Iterator()
		for iterator.HasNext() {
			id, _ := iterator.Next()
			_, has := cards[id]
			if has {
				cards[id] = 1
			} else {
				cards[id]++
			}
		}
	}
	stacks = []list.List[int32]{deck.Main, deck.Extra, deck.Side}
	for _, stack := range stacks {
		iterator := stack.Iterator()
		for iterator.HasNext() {
			id, _ := iterator.Next()
			_, has := ncards[id]
			if has {
				ncards[id] = 1
			} else {
				ncards[id]++
			}
		}
	}
	for k, v := range cards {
		if _, has := ncards[k]; !has {
			return false
		}
		if ncards[k] != v {
			return false
		}

	}
	return true
}

func (d *Deck) AddToCards(cards map[int32]int, card *Card) {
	id := card.Id
	if card.Alias != 0 {
		id = card.Alias
	}
	if _, has := cards[id]; has {
		cards[id]++
	} else {
		cards[id] = 1
	}
}
