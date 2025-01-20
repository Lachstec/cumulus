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

type PvP string

const (
	Enabled  PvP = "true"  //nolint:all
	Disabled PvP = "false" //nolint:all
)

type Image string

const (
	Alpine3_20_3 Image = "29a24dc0-b24b-4cc8-b43b-a8a4c6916d0f" //nolint:all
	Ubuntu14_04  Image = "f211d66d-9167-4133-abd6-c40b1586394e" //nolint:all
	Ubuntu16_04  Image = "cc21acd3-8cf1-40be-8bb1-6fea9453b0bb" //nolint:all
	Ubuntu18_04  Image = "f56402aa-369e-42b7-a64d-2db29bebfebd" //nolint:all
	Ubuntu20_04  Image = "de55d4f6-5905-4e4e-b42b-db7ebabbcda4" //nolint:all
	Ubuntu22_04  Image = "1404d277-1fd2-4682-9fbd-80c7d05b80e1" //nolint:all
	Ubuntu24_04  Image = "d6d1835c-7180-4ca9-b4a1-470afbd8b398" //nolint:all
	Fedora41     Image = "715efd2c-b224-49d6-bbc7-00c204b1f04c" //nolint:all
	Cirros063    Image = "b3f88062-56d5-43b1-9d1e-96980ea0e16b" //nolint:all
)

type Server struct {
	ID               int64
	UserID           int64        `db:"user_id"`
	OpenstackID      string       `db:"openstack_id"`
	Name             string       `db:"name" json:"name"`
	Address          net.IP       `db:"addr"`
	Status           ServerStatus `db:"server_status"`
	Port             int          `db:"port"`
	Flavour          Flavour
	Image            Image      `db:"image" json:"image"`
	Game             string     `db:"game" json:"game"`
	GameVersion      string     `db:"game_version" json:"game_version"`
	GameMode         GameMode   `db:"game_mode" json:"gamemode"`
	Difficulty       Difficulty `db:"difficulty" json:"difficulty"`
	WhitelistEnabled bool       `db:"whitelist_enabled" json:"whitelist_enabled"`
	PvPEnabled       PvP        `db:"pvp_enabled" json:"pvp_enabled"`
	PlayersMax       int        `db:"players_max" json:"players_max"`
	SSHKey           []byte     `db:"ssh_key"`
}
