package pvp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"val-snitch/internal/constants"
)

func getCurrentSeasonID(info constants.LogInfo, accessToken string, entitlement string) string {
	url := fmt.Sprintf("https://shared.%s.a.pvp.net/content-service/v3/content", info.Shard)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req.Header.Add("X-Riot-Entitlements-JWT", entitlement)
	req.Header.Add("X-Riot-ClientVersion", info.ClientVersion)
	req.Header.Add("X-Riot-ClientPlatform", constants.DefaultClientPlatformString)
	req.Header.Add("User-Agent", "")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var seasons map[string]interface{}
	if err = json.Unmarshal(body, &seasons); err != nil {
		return ""
	}

	if seasonsSlice, ok := seasons["Seasons"].([]interface{}); ok {
		for _, elem := range seasonsSlice {
			if season, ok := elem.(map[string]interface{}); ok {
				if season["Type"] == "act" && season["IsActive"] == true {
					return season["ID"].(string)
				}
			}
		}
	}
	return ""
}
