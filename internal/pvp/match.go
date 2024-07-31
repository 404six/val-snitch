package pvp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"val-snitch/internal/constants"
)

func getCurrentMatchID(info constants.LogInfo, accessToken string, entitlement string) (string, error) {
	url := fmt.Sprintf("https://glz-%s-1.%s.a.pvp.net/core-game/v1/players/%s", info.Region, info.Shard, info.Puuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("X-Riot-Entitlements-JWT", entitlement)
	req.Header.Add("X-Riot-ClientVersion", info.ClientVersion)
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

func GetCurrentGameMatch(info constants.LogInfo, accessToken string, entitlement string) constants.MatchInfo {
	mi := constants.MatchInfo{}

	currentMatchID, err := getCurrentMatchID(info, accessToken, entitlement)
	if err != nil {
		return mi
	}

	url := fmt.Sprintf("https://glz-%s-1.%s.a.pvp.net/core-game/v1/matches/%s", info.Region, info.Shard, currentMatchID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return mi
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("X-Riot-Entitlements-JWT", entitlement)
	req.Header.Add("X-Riot-ClientVersion", info.ClientVersion)
	req.Header.Add("X-Riot-ClientPlatform", constants.DefaultClientPlatformString)
	req.Header.Add("User-Agent", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return mi
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return mi
	}

	var match map[string]interface{}
	if err = json.Unmarshal(body, &match); err != nil {
		return mi
	}

	mi.ID = match["MatchID"].(string)
	mi.MapID = match["MapID"].(string)
	mi.ModeID = match["ModeID"].(string)

	if playersSlice, ok := match["Players"].([]interface{}); ok {
		for _, elem := range playersSlice {
			if player, ok := elem.(map[string]interface{}); ok {
				playerInfo := constants.PlayerInfo{
					Puuid:  player["Subject"].(string),
					Agent:  player["CharacterID"].(string),
					TeamID: player["TeamID"].(string),
					Rank:   GetPlayerRankByUUID(info, accessToken, entitlement, player["Subject"].(string)),
				}

				// get account level
				for k, v := range player["PlayerIdentity"].(map[string]interface{}) {
					if k == "AccountLevel" {
						playerInfo.AccountLevel = v.(float64)
						break
					}
				}

				mi.Players = append(mi.Players, playerInfo)
			}
		}
	}

	return mi
}
