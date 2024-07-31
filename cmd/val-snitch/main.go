package main

import (
	"fmt"
	"val-snitch/internal/auth"
	"val-snitch/internal/pvp"
	"val-snitch/internal/utils"
)

func main() {
	logInfo := auth.GetClientInfo()
	fmt.Printf("logInfo: %+v\n", logInfo)

	accessToken, err := auth.AuthFromClient()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	entitlementToken, err := auth.GetEntitlement(accessToken)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	mi := pvp.GetCurrentGameMatch(logInfo, accessToken, entitlementToken)

	// playerInfo
	for _, elem := range mi.Players {
		fmt.Printf("puuid: %v | rank: %v | agent: %v | rankName: %v | accountLevel: %v | teamID: %v\n",
			elem.Puuid, elem.Rank, elem.Agent, utils.GetRankName(int(elem.Rank)), elem.AccountLevel, elem.TeamID)
	}
}
