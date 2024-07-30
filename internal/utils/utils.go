package utils

import (
	"os"
	"path/filepath"
)

func Get_log_path() string {
	local_appdata := os.Getenv("LOCALAPPDATA")

	file_path := filepath.Join(local_appdata, "VALORANT", "Saved", "Logs", "ShooterGame.log")

	return file_path
}
