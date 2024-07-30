package constants

type Player_info struct {
	Puuid         string  `json:"puuid"`
	Rank          float64 `json:"rank"`
	Agent         string  `json:"agent"`
	Account_level float64 `json:"account_level"`
	Team_id       string  `json:"team_id"`
}
