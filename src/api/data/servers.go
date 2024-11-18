package data

type Server struct {
	ID int
	Name string `json:"name"`
	IP string `json:"ip"`	
	Version string `json:"version"`
	Mode string `json:"mode"`
	Difficulty string `json:"difficulty"`
	MaxPlayers int `json:"maxplayers"`
	PvP bool `json:"pvp"`
}

var Servers []Server
	
func init() {
	Servers = []Server{
		{
			ID: 1,
			Name: "Craftopia",
			IP: "192.168.1.1",
			Version: "1.20.1",
			Mode: "Survival",
			Difficulty: "Hard",
			MaxPlayers: 20,
			PvP: true,
		  },
		  {
			ID: 2,
			Name: "Skyland Adventures",
			IP: "203.0.113.5",
			Version: "1.19.3",
			Mode: "Adventure",
			Difficulty: "Normal",
			MaxPlayers: 15,
			PvP: false,
		  },
		  {
			ID: 3,
			Name: "PixelVerse",
			IP: "198.51.100.23",
			Version: "1.18.2",
			Mode: "Creative",
			Difficulty: "Peaceful",
			MaxPlayers: 25,
			PvP: false,
		  },
		  {
			ID: 4,
			Name: "EpicWars",
			IP: "203.0.113.77",
			Version: "1.20.1",
			Mode: "Survival",
			Difficulty: "Normal",
			MaxPlayers: 30,
			PvP: true,
		  },
		  {
			ID: 5,
			Name: "MysticCraft",
			IP: "198.51.100.11",
			Version: "1.19",
			Mode: "Adventure",
			Difficulty: "Easy",
			MaxPlayers: 10,
			PvP: false,
		  },
		  {
			ID: 6,
			Name: "NetherQuest",
			IP: "192.168.0.45",
			Version: "1.20",
			Mode: "Survival",
			Difficulty: "Hard",
			MaxPlayers: 5,
			PvP: true,
		  },
		  {
			ID: 7,
			Name: "Blocky Realms",
			IP: "198.51.100.55",
			Version: "1.18.1",
			Mode: "Creative",
			Difficulty: "Normal",
			MaxPlayers: 15,
			PvP: false,
		  },
		  {
			ID: 8,
			Name: "BedWars Arena",
			IP: "203.0.113.90",
			Version: "1.20.1",
			Mode: "Adventure",
			Difficulty: "Normal",
			MaxPlayers: 25,
			PvP: true,
		  },
		  {
			ID: 9,
			Name: "VillageCraft",
			IP: "192.0.2.60",
			Version: "1.19.2",
			Mode: "Survival",
			Difficulty: "Peaceful",
			MaxPlayers: 12,
			PvP: false,
		  },
		  {
			ID: 10,
			Name: "Adventure Realm",
			IP: "198.51.100.200",
			Version: "1.18.2",
			Mode: "Adventure",
			Difficulty: "Hard",
			MaxPlayers: 20,
			PvP: true,
		  },
	}
}