package constants

import "regexp"

type LogInfo struct {
	Puuid          string
	Shard          string
	Region         string
	Client_version string
}

const DefaultClientPlatformString = "ew0KCSJwbGF0Zm9ybVR5cGUiOiAiUEMiLA0KCSJwbGF0Zm9ybU9TIjogIldpbmRvd3MiLA0KCSJwbGF0Zm9ybU9TVmVyc2lvbiI6ICIxMC4wLjE5MDQyLjEuMjU2LjY0Yml0IiwNCgkicGxhdGZvcm1DaGlwc2V0IjogIlVua25vd24iDQp9"

var Ssid_regex = regexp.MustCompile(`name:\s*\"ssid\"\s*[\s\S]*?value:\s*\"([^\"]+)\"`)

var (
	Puuid_regex          = regexp.MustCompile(`Logged in user changed: (.+)`)
	Region_shard_regex   = regexp.MustCompile(`https://glz-(.+?)-1\.(.+?)\.a\.pvp\.net`)
	Client_version_regex = regexp.MustCompile(`CI server version: (.+)`)
)
