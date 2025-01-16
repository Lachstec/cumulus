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

type Flavour struct {
	ID string
	Name string
	RAM	int
	Disk int
	Ephemeral int
	VCPUs int
	Is_Public bool
}

var Flavours = []Flavour {
	{
		ID: "1",
		Name: "m1.tiny",
		RAM: 512,
		Disk: 1,
		Ephemeral: 0,
		VCPUs: 1,
		Is_Public: true,
	},
	{
		ID: "2",
		Name: "m1.small",
		RAM: 2048,
		Disk: 20,
		Ephemeral: 0,
		VCPUs: 1,
		Is_Public: true,
	},
	{
		ID: "3",
		Name: "m1.medium",
		RAM: 4096,
		Disk: 40,
		Ephemeral: 0,
		VCPUs: 2,
		Is_Public: true,
	},
	{
		ID: "4",
		Name: "m1.large",
		RAM: 8192,
		Disk: 80,
		Ephemeral: 0,
		VCPUs: 4,
		Is_Public: true,
	},
	{
		ID: "5",
		Name: "m1.xlarge",
		RAM: 16384,
		Disk: 160,
		Ephemeral: 0,
		VCPUs: 8,
		Is_Public: true,
	},
}

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
	OpenstackID      string `db:"openstack_id"`
	Name             string
	Address          net.IP `db:"addr"`
	Status           ServerStatus
	Port             int
	Flavour			 Flavour
	Image			 Image
	Memory           int `db:"memory_mb"`
	Game             string
	GameVersion      string   `db:"game_version"`
	GameMode         GameMode `db:"game_mode"`
	Difficulty       Difficulty
	WhitelistEnabled bool `db:"whitelist_enabled"`
	PlayersMax       int  `db:"players_max"`
	PvP              PvP
	SSHKey           []byte `db:"ssh_key"`
}
