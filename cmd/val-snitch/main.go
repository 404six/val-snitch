package main

import (
	"fmt"
	"val-snitch/internal/auth"
	"val-snitch/internal/pvp"
	"val-snitch/internal/utils"
)

func main() {

	log_info := auth.Get_client_info()
	fmt.Printf("log_info: %+v\n", log_info)

	access_token, err := auth.Auth_from_client()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	entitlement_token, err := auth.Get_entitlement(access_token)

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	mi := pvp.Get_current_game_match(log_info, access_token, entitlement_token)

	// player_info
	for _, elem := range mi.Players {
		fmt.Printf("ppuid: %v | rank: %v | agent: %v | rank_name: %v | account_level: %v | team_id: %v\n", elem.Puuid, elem.Rank, elem.Agent, utils.Get_rank_name(int(elem.Rank)), elem.Account_level, elem.Team_id)
	}

}
