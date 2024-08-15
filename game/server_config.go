package game

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type ServerConfig struct {
	ServerPort    int
	Path          string
	ScriptFolder  string
	CardCDB       string
	BanlistFile   string
	Log           bool
	ConsoleLog    bool
	HandShuffle   bool
	AutoEndTurn   bool
	ClientVersion int
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		ClientVersion: 0x1332,
		ServerPort:    8911,
		Path:          ".",
		ScriptFolder:  "script",
		CardCDB:       "cards.cdb",
		BanlistFile:   "lflist.conf",
		Log:           true,
		ConsoleLog:    true,
		HandShuffle:   false,
		AutoEndTurn:   true,
	}
}

func (config *ServerConfig) Load(file string) bool {
	if file == "" {
		file = "config.txt"
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}

	fileHandle, err := os.Open(file)
	if err != nil {
		Logger.WriteError(err)
		return false
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		data := strings.SplitN(line, "=", 2)
		variable := strings.ToLower(strings.TrimSpace(data[0]))
		value := strings.TrimSpace(data[1])

		switch variable {
		case "serverport":
			if port, err := strconv.Atoi(value); err == nil {
				config.ServerPort = port
			}
		case "path":
			config.Path = value
		case "scriptfolder":
			config.ScriptFolder = value
		case "cardcdb":
			config.CardCDB = value
		case "banlist":
			config.BanlistFile = value
		case "errorlog":
			if logValue, err := strconv.ParseBool(value); err == nil {
				config.Log = logValue
			}
		case "consolelog":
			if consoleLogValue, err := strconv.ParseBool(value); err == nil {
				config.ConsoleLog = consoleLogValue
			}
		case "handshuffle":
			if handShuffleValue, err := strconv.ParseBool(value); err == nil {
				config.HandShuffle = handShuffleValue
			}
		case "autoendturn":
			if autoEndTurnValue, err := strconv.ParseBool(value); err == nil {
				config.AutoEndTurn = autoEndTurnValue
			}
		case "clientversion":
			if clientVersion, err := strconv.ParseInt(value, 16, 0); err == nil {
				config.ClientVersion = int(clientVersion)
			}
		}
	}

	if config.HandShuffle {
		Logger.WriteLine("Warning: Hand shuffle requires a custom ocgcore to work.")
	}
	return true
}
