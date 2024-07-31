package constants

type MatchInfo struct {
	ID      string       `json:"id"`
	MapID   string       `json:"map_id"`
	ModeID  string       `json:"mode_id"`
	Players []PlayerInfo `json:"players"`
}

var Ranks = []string{
	"unranked",
	"iron",
	"bronze",
	"silver",
	"gold",
	"platinum",
	"diamond",
	"ascendent",
	"immortal",
	"radiant",
}
