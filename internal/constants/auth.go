package constants

import "regexp"

type LogInfo struct {
	Puuid          string
	Shard          string
	Region         string
	Client_version string
}

var (
	Puuid_regex          = regexp.MustCompile(`Logged in user changed: (.+)`)
	Region_shard_regex   = regexp.MustCompile(`https://glz-(.+?)-1\.(.+?)\.a\.pvp\.net`)
	Client_version_regex = regexp.MustCompile(`CI server version: (.+)`)
)
