package utils

import (
	"os"
	"path/filepath"
	"strings"
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
