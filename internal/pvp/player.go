package pvp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"val-snitch/internal/constants"
)

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

	return mi
}
