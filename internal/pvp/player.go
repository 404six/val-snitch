package pvp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"val-snitch/internal/constants"
)

func Get_player_rank_by_uuid(info constants.LogInfo, accessToken string, entitlement string, puuid string) float64 {

	url := fmt.Sprintf("https://pd.%s.a.pvp.net/mmr/v1/players/%s", info.Shard, puuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("X-Riot-Entitlements-JWT", entitlement)
	req.Header.Add("X-Riot-ClientVersion", info.Client_version)
	req.Header.Add("X-Riot-ClientPlatform", constants.DefaultClientPlatformString)
	req.Header.Add("User-Agent", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return -1
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1
	}

	season_id := get_current_season_id(info, accessToken, entitlement)

	var queue_skills map[string]interface{}
	if err = json.Unmarshal(body, &queue_skills); err != nil {
		return -1
	}

	if qs, ok := queue_skills["QueueSkills"].(map[string]interface{}); ok {
		if competitive, ok := qs["competitive"].(map[string]interface{}); ok {
			if season_slice, ok := competitive["SeasonalInfoBySeasonID"].(map[string]interface{}); ok {
				for _, elem := range season_slice {
					if season, ok := elem.(map[string]interface{}); ok {
						if season["SeasonID"] == season_id {
							return season["Rank"].(float64)
						}
					}
				}
			}
		}
	}
	// competitive not unlocked yet
	return 999
}
