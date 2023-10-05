package duel

import (
	"bufio"
	"fmt"
	"github.com/spf13/cast"
	"os"
)

type LFList struct {
	Hash     uint32
	ListName string
	Content  map[uint32]uint32
}

type DeckManager struct {
	DeckBuffer [65536]byte
	LFList     []LFList
}

func (dm *DeckManager) LoadLFListSingle(path string) {
	var cur *LFList
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || line[0] == '#' {
			continue
		}

		if line[0] == '!' {
			strBuffer := make([]rune, 256)
			sa := decodeUTF8(line[1:], strBuffer)
			for sa > 0 && (strBuffer[sa-1] == '\r' || strBuffer[sa-1] == '\n') {
				sa--
			}
			strBuffer = strBuffer[:sa]
			newList := LFList{}
			dm.LFList = append(dm.LFList, newList)
			cur = &dm.LFList[len(dm.LFList)-1]
			cur.ListName = string(strBuffer)
			cur.Hash = 0x7dfcee6a
			continue
		}

		p := 0
		for p < len(line) && line[p] != ' ' && line[p] != '\t' {
			p++
		}
		if p >= len(line) {
			continue
		}
		lineBuf := line[:p]
		p++
		sa := p
		code, err := cast.ToUint32E(lineBuf)
		if err != nil || code == 0 {
			continue
		}
		for p < len(line) && (line[p] == ' ' || line[p] == '\t') {
			p++
		}
		for p < len(line) && line[p] != ' ' && line[p] != '\t' {
			p++
		}
		count, err := cast.ToUint32E(line[sa:p])
		if err != nil {
			continue
		}
		if cur == nil {
			continue
		}
		cur.Content[code] = count
		cur.Hash = cur.Hash ^ ((code << 18) | (code >> 14)) ^ ((code << (27 + count)) | (code >> (5 - count)))
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func decodeUTF8(input string, output []rune) int {
	copy(output, []rune(input))
	return len(input)
}
