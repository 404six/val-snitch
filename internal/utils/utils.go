package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"val-snitch/internal/constants"
)

func Get_log_path() string {
	local_appdata := os.Getenv("LOCALAPPDATA")

	file_path := filepath.Join(local_appdata, "VALORANT", "Saved", "Logs", "ShooterGame.log")

	return file_path
}

func Get_settings_path() string {
	local_appdata := os.Getenv("LOCALAPPDATA")

	file_path := filepath.Join(local_appdata, "Riot Games", "Riot Client", "Data", "RiotGamesPrivateSettings.yaml")

	return file_path
}

func Ellipsis_str(str string, length int) string {
	if len(str) > length {
		return str[:length-3] + "..."
	}
	return str
}

func Get_string_between(str, start_delim, end_delim string) string {
	start_index := strings.Index(str, start_delim)
	start_index += len(start_delim)

	end_index := strings.Index(str[start_index:], end_delim)
	end_index += start_index
	return str[start_index:end_index]
}

func Get_rank_name(index int) string {
	if index == 0 {
		return constants.Ranks[0]
	} else if index == 27 {
		return constants.Ranks[9]
	}
	base_rank_index := (index - 3) / 3
	level := (index-3)%3 + 1
	return fmt.Sprintf("%s %v", constants.Ranks[base_rank_index+1], level)
}
