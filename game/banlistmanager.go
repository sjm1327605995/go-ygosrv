package game

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type BanlistManager struct {
	Banlists []Banlist
}

var BanListManager *BanlistManager

func InitBanListManager(fileName string) {
	BanListManager = new(BanlistManager)
	BanListManager.Banlists = []Banlist{}
	var current *Banlist
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "!") {
			current = &Banlist{}
			BanListManager.Banlists = append(BanListManager.Banlists, *current)
			continue
		}
		if !strings.Contains(line, " ") || current == nil {
			continue
		}
		data := strings.Fields(line)
		id := parseInt(data[0])
		count := parseInt(data[1])
		current.Add(id, count)
	}
}

// GetIndex returns the index of the Banlist that matches the specified hash
// <ai>Preserve original comments if possible</ai>
func (bm *BanlistManager) GetIndex(hash uint32) int {
	for i := 0; i < len(bm.Banlists); i++ {
		if bm.Banlists[i].Hash == hash {
			return i
		}
	}
	return 0
}

// Helper function to parse an integer from a string safely
func parseInt(value string) int {
	result, _ := strconv.Atoi(value)
	return result
}
