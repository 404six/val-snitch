package pvp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"val-snitch/internal/constants"
)

func GetPlayerRankByUUID(info constants.LogInfo, accessToken string, entitlement string, puuid string) float64 {
	url := fmt.Sprintf("https://pd.%s.a.pvp.net/mmr/v1/players/%s", info.Shard, puuid)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return -1
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("X-Riot-Entitlements-JWT", entitlement)
	req.Header.Add("X-Riot-ClientVersion", info.ClientVersion)
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

	seasonID := getCurrentSeasonID(info, accessToken, entitlement)

	var queueSkills map[string]interface{}
	if err = json.Unmarshal(body, &queueSkills); err != nil {
		return -1
	}

	if qs, ok := queueSkills["QueueSkills"].(map[string]interface{}); ok {
		if competitive, ok := qs["competitive"].(map[string]interface{}); ok {
			if seasonalInfo, ok := competitive["SeasonalInfoBySeasonID"].(map[string]interface{}); ok {
				for _, elem := range seasonalInfo {
					if season, ok := elem.(map[string]interface{}); ok {
						if season["SeasonID"] == seasonID {
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
