package game

import (
	"strconv"
	"strings"
)

type GameConfig struct {
	LfList         int
	Rule           int
	Mode           int
	EnablePriority bool
	NoCheckDeck    bool
	NoShuffleDeck  bool
	StartLp        int
	StartHand      int
	DrawCount      int
	GameTimer      int
	Name           string
}

func NewGameConfig(info string) *GameConfig {
	config := &GameConfig{}

	if strings.ToLower(info) == "tcg" || strings.ToLower(info) == "ocg" || strings.ToLower(info) == "ocg/tcg" || strings.ToLower(info) == "tcg/ocg" {
		if strings.ToLower(info) == "ocg/tcg" {
			config.LfList = 1
		} else if strings.ToLower(info) == "tcg/ocg" {
			config.LfList = 0
		} else {
			config.LfList = 0
		}
		if strings.ToLower(info) == "ocg/tcg" || strings.ToLower(info) == "tcg/ocg" {
			config.Rule = 2
		} else if strings.ToLower(info) == "tcg" {
			config.Rule = 1
		} else {
			config.Rule = 0
		}
		config.Mode = 0
		config.EnablePriority = false
		config.NoCheckDeck = false
		config.NoShuffleDeck = false
		config.StartLp = 8000
		config.StartHand = 5
		config.DrawCount = 1
		config.GameTimer = 120
		config.Name = &GameManager{}.RandomRoomName()
	} else {
		config.Load(info)
	}

	return config
}

func NewGameConfigFromPacket(packet *GameClientPacket) *GameConfig {
	config := &GameConfig{}
	config.LfList = BanlistManager.GetIndex(packet.ReadUInt32())
	config.Rule = int(packet.ReadByte())
	config.Mode = int(packet.ReadByte())
	config.EnablePriority = packet.ReadByte() != 0
	config.NoCheckDeck = packet.ReadByte() != 0
	config.NoShuffleDeck = packet.ReadByte() != 0
	// C++ padding: 5 bytes + 3 bytes = 8 bytes
	for i := 0; i < 3; i++ {
		packet.ReadByte()
	}
	config.StartLp = int(packet.ReadInt32())
	config.StartHand = int(packet.ReadByte())
	config.DrawCount = int(packet.ReadByte())
	config.GameTimer = int(packet.ReadInt16())
	packet.ReadUnicode(20)
	config.Name = packet.ReadUnicode(30)
	if config.Name == "" {
		config.Name = GameManager.RandomRoomName()
	}
	return config
}

func (config *GameConfig) Load(gameInfo string) {
	defer func() {
		if r := recover(); r != nil {
			config.LfList = 0
			config.Rule = 2
			config.Mode = 0
			config.EnablePriority = false
			config.NoCheckDeck = false
			config.NoShuffleDeck = false
			config.StartLp = 8000
			config.StartHand = 5
			config.DrawCount = 1
			config.GameTimer = 120
			config.Name = GameManager.RandomRoomName()
		}
	}()

	rules := gameInfo[:6]

	config.Rule, _ = strconv.Atoi(string(rules[0]))
	config.Mode, _ = strconv.Atoi(string(rules[1]))
	if rules[2] == '0' {
		config.GameTimer = 120
	} else {
		config.GameTimer = 60
	}
	config.EnablePriority = rules[3] == 'T' || rules[3] == '1'
	config.NoCheckDeck = rules[4] == 'T' || rules[4] == '1'
	config.NoShuffleDeck = rules[5] == 'T' || rules[5] == '1'

	data := gameInfo[6:]
	list := strings.Split(data, ",")

	config.StartLp, _ = strconv.Atoi(list[0])
	config.LfList, _ = strconv.Atoi(list[1])
	config.StartHand, _ = strconv.Atoi(list[2])
	config.DrawCount, _ = strconv.Atoi(list[3])
	config.Name = list[4]
}
