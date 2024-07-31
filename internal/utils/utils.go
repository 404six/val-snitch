package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"val-snitch/internal/constants"
)

func GetLogPath() string {
	localAppData := os.Getenv("LOCALAPPDATA")

	filePath := filepath.Join(localAppData, "VALORANT", "Saved", "Logs", "ShooterGame.log")

	return filePath
}

func GetSettingsPath() string {
	localAppData := os.Getenv("LOCALAPPDATA")

	filePath := filepath.Join(localAppData, "Riot Games", "Riot Client", "Data", "RiotGamesPrivateSettings.yaml")

	return filePath
}

func EllipsisStr(str string, length int) string {
	if len(str) > length {
		return str[:length-3] + "..."
	}
	return str
}

func GetStringBetween(str, startDelim, endDelim string) string {
	startIndex := strings.Index(str, startDelim)
	startIndex += len(startDelim)

	endIndex := strings.Index(str[startIndex:], endDelim)
	endIndex += startIndex
	return str[startIndex:endIndex]
}

func GetRankName(index int) string {
	if index == 0 {
		return constants.Ranks[0]
	} else if index == 27 {
		return constants.Ranks[9]
	}
	baseRankIndex := (index - 3) / 3
	level := (index-3)%3 + 1
	return fmt.Sprintf("%s %v", constants.Ranks[baseRankIndex+1], level)
}
