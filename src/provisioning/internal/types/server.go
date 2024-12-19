package types

import "net"

type ServerStatus string

const (
	Running    ServerStatus = "running"    //nolint:all
	Stopped    ServerStatus = "stopped"    //nolint:all
	Restarting ServerStatus = "restarting" //nolint:all
)

type GameMode string

const (
	Creative  GameMode = "creative" //nolint:all
	Survival  GameMode = "survival" //nolint:all
	Adventure GameMode = "peaceful" //nolint:all
	Hardcore  GameMode = "hardcore" //nolint:all
)

type Difficulty string

const (
	Peaceful Difficulty = "peaceful" //nolint:all
	Easy     Difficulty = "easy"     //nolint:all
	Normal   Difficulty = "normal"   //nolint:all
	Hard     Difficulty = "hard"     //nolint:all
)

type Server struct {
	Id               int64
	OpenstackId      string `db:"openstack_id"`
	Name             string
	Address          net.IP `db:"addr"`
	Status           ServerStatus
	Port             int
	Memory           int `db:"memory_mb"`
	Game             string
	GameVersion      string   `db:"game_version"`
	GameMode         GameMode `db:"game_mode"`
	Difficulty       Difficulty
	WhitelistEnabled bool   `db:"whitelist_enabled"`
	PlayersMax       int    `db:"players_max"`
	SSHKey           []byte `db:"ssh_key"`
}
