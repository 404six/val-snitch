package pvp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"val-snitch/internal/constants"
)

func get_current_match_id(info constants.LogInfo, access_token string, entitlement string) (string, error) {

	url := fmt.Sprintf("https://glz-%s-1.%s.a.pvp.net/core-game/v1/players/%s", info.Region, info.Shard, info.Puuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+access_token)
	req.Header.Add("X-Riot-Entitlements-JWT", entitlement)
	req.Header.Add("X-Riot-ClientVersion", info.Client_version)
	req.Header.Add("X-Riot-ClientPlatform", constants.DefaultClientPlatformString)
	req.Header.Add("User-Agent", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("player is not in an active match (after the agent select screen)")
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get current game match ID: %d %s - %s", resp.StatusCode, resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	matchID, ok := responseData["MatchID"].(string)
	if !ok {
		return "", errors.New("MatchID not found in response")
	}

	return matchID, nil
}
