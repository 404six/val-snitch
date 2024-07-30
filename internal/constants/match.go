package constants

type Match_info struct {
	Id      string        `json:"id"`
	Map_id  string        `json:"map_id"`
	Mode_id string        `json:"mode_id"`
	Players []Player_info `json:"players"`
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
