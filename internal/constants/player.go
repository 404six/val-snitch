package constants

type PlayerInfo struct {
	Puuid        string  `json:"puuid"`
	Rank         float64 `json:"rank"`
	Agent        string  `json:"agent"`
	AccountLevel float64 `json:"account_level"`
	TeamID       string  `json:"team_id"`
}
