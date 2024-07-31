package auth

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"val-snitch/internal/constants"
	"val-snitch/internal/utils"
)

func GetEntitlement(accessToken string) (string, error) {
	url := "https://entitlements.auth.riotgames.com/api/token/v1"
	body := []byte("{}")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	entitlementsToken := utils.GetStringBetween(string(respBody), `"entitlements_token":"`, `"}`)

	return entitlementsToken, nil
}

func ssidReauth(ssid string) (string, error) {
	req, err := http.NewRequest("GET", "https://auth.riotgames.com/authorize?redirect_uri=https%3A%2F%2Fplayvalorant.com%2Fopt_in&client_id=play-valorant-web-prod&response_type=token%20id_token&nonce=1&scope=account%20openid", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "")
	req.Header.Set("Cookie", "ssid="+ssid)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" {
		return "", errors.New("no location header in response")
	}
	if !strings.HasPrefix(location, "https://playvalorant.com/opt_in") {
		return "", fmt.Errorf("invalid reauth location: %s", utils.EllipsisStr(location, 40))
	}
	return location, nil
}

func AuthFromClient() (string, error) {
	settingsPath := utils.GetSettingsPath()
	settings, err := os.ReadFile(settingsPath)
	if err != nil {
		return "", fmt.Errorf("failed to read settings file: %v", err)
	}

	match := constants.SsidRegex.FindSubmatch(settings)
	if match == nil {
		return "", errors.New("couldn't find ssid in RiotGamesPrivateSettings.yaml")
	}

	ssid := string(match[1])
	if len(ssid) == 0 || len(strings.Split(ssid, ".")) != 3 {
		return "", fmt.Errorf("invalid ssid: %s", ssid)
	}

	// From https://github.com/techchrism/riot-auth-test; the ssid reauth might fail but works on a retry
	var errors []error
	for i := 0; i < 3; i++ {
		result, err := ssidReauth(ssid)
		if err == nil {
			accessToken := utils.GetStringBetween(result, "access_token=", "&scope")
			return accessToken, nil
		}
		errors = append(errors, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return "", fmt.Errorf("failed to reauth after %d attempts: %v", len(errors), errors)
}

func GetClientInfo() constants.LogInfo {
	logInfo := constants.LogInfo{}

	filePath := utils.GetLogPath()
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return logInfo
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if match := constants.PuuidRegex.FindStringSubmatch(line); match != nil {
			logInfo.Puuid = match[1]
		}

		if match := constants.RegionShardRegex.FindStringSubmatch(line); match != nil {
			logInfo.Region = match[1]
			logInfo.Shard = match[2]
		}

		if match := constants.ClientVersionRegex.FindStringSubmatch(line); match != nil {
			re := regexp.MustCompile(`^(release-\d+\.\d+-)`)
			logInfo.ClientVersion = re.ReplaceAllString(match[1], "${1}shipping-")
		}
	}
	return logInfo
}
