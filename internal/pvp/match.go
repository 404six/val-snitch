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

func Get_current_game_match(info constants.LogInfo, access_token string, entitlement string) constants.Match_info {

	mi := constants.Match_info{}

	current_match_id, err := get_current_match_id(info, access_token, entitlement)
	if err != nil {
		fmt.Printf("error 0 - %v", fmt.Errorf("%v", err))
		return mi
	}

	url := fmt.Sprintf("https://glz-%s-1.%s.a.pvp.net/core-game/v1/matches/%s", info.Region, info.Shard, current_match_id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("error 1 - %v", fmt.Errorf("%v", err))
		return mi
	}
	req.Header.Add("Authorization", "Bearer "+access_token)
	req.Header.Add("X-Riot-Entitlements-JWT", entitlement)
	req.Header.Add("X-Riot-ClientVersion", info.Client_version)
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

	mi.Id = match["MatchID"].(string)
	mi.Map_id = match["MapID"].(string)
	mi.Mode_id = match["ModeID"].(string)

	if players_slice, ok := match["Players"].([]interface{}); ok {
		for _, elem := range players_slice {
			if player, ok := elem.(map[string]interface{}); ok {

				player_info := constants.Player_info{
					Puuid:   player["Subject"].(string),
					Agent:   player["CharacterID"].(string),
					Team_id: player["TeamID"].(string),
					Rank:    Get_player_rank_by_uuid(info, access_token, entitlement, player["Subject"].(string)),
				}

				// get account level
				for k, v := range player["PlayerIdentity"].(map[string]interface{}) {
					if k == "AccountLevel" {
						player_info.Account_level = v.(float64)
						break
					}
				}

				mi.Players = append(mi.Players, player_info)
			}
		}
	}

	return mi
}
