package auth

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"val-snitch/internal/constants"
	"val-snitch/internal/utils"
)

func Get_client_info() constants.LogInfo {
	log_info := constants.LogInfo{}

	file_path := utils.Get_log_path()
	file, err := os.Open(file_path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return log_info
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if match := constants.Puuid_regex.FindStringSubmatch(line); match != nil {
			log_info.Puuid = match[1]
		}

		if match := constants.Region_shard_regex.FindStringSubmatch(line); match != nil {
			log_info.Region = match[1]
			log_info.Shard = match[2]
		}

		if match := constants.Client_version_regex.FindStringSubmatch(line); match != nil {
			re := regexp.MustCompile(`^(release-\d+\.\d+-)`)
			log_info.Client_version = re.ReplaceAllString(match[1], "${1}shipping-")
		}
	}
	return log_info
}
