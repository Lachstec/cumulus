package types

import "net"

type ServerStatus string

const (
	running    ServerStatus = "running"
	stopped    ServerStatus = "stopped"
	restarting ServerStatus = "restarting"
)

type GameMode string

const (
	creative  GameMode = "creative"
	survival  GameMode = "survival"
	adventure GameMode = "peaceful"
	hardcore  GameMode = "hardcore"
)

type Difficulty string

const (
	peaceful Difficulty = "peaceful"
	easy     Difficulty = "easy"
	normal   Difficulty = "normal"
	hard     Difficulty = "hard"
)

type Server struct {
	Id               int64
	Name             string
	Address          net.IPNet
	Status           ServerStatus
	Port             int
	Memory           int
	Game             string
	GameVersion      string
	GameMode         GameMode
	Difficulty       Difficulty
	WhitelistEnabled bool
	PlayersMax       int
}
