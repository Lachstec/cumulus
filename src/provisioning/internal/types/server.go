package types

import "net"

type ServerStatus string

const (
	Running    ServerStatus = "running"
	Stopped    ServerStatus = "stopped"
	Restarting ServerStatus = "restarting"
)

type GameMode string

const (
	Creative  GameMode = "creative"
	Survival  GameMode = "survival"
	Adventure GameMode = "peaceful"
	Hardcore  GameMode = "hardcore"
)

type Difficulty string

const (
	Peaceful Difficulty = "peaceful"
	Easy     Difficulty = "easy"
	Normal   Difficulty = "normal"
	Hard     Difficulty = "hard"
)

type Server struct {
	Id               int64
	Name             string
	Address          net.IP `db:"addr"`
	Status           ServerStatus
	Port             int
	Memory           int `db:"memory_mb"`
	Game             string
	GameVersion      string   `db:"game_version"`
	GameMode         GameMode `db:"game_mode"`
	Difficulty       Difficulty
	WhitelistEnabled bool `db:"whitelist_enabled"`
	PlayersMax       int  `db:"players_max"`
}
