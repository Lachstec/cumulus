package services

import (
	Data "data"
)

func ReadNumOfServers() int {
	return len(Data.Servers)
}

func ReadAllServers() []Data.Server {
	return Data.Servers
}

func ReadServerByServerID(serverid int) Data.Server {
	return Data.Servers[serverid]
}

func CreateServer(server Data.Server) {
	Data.Servers = append(Data.Servers, server)
}

func DeleteServerByServerID(serverid int) {
	Data.Servers = append(Data.Servers[:serverid], Data.Servers[serverid + 1:]...)
}

func UpdateServer(serverid int, server Data.Server) {
	switch {
		case server.Name != "":
			Data.Servers[serverid].Name = server.Name
		case server.Difficulty != "":
			Data.Servers[serverid].Difficulty = server.Difficulty
		case server.IP != "":
			Data.Servers[serverid].IP = server.IP
		case server.MaxPlayers != 0:
			Data.Servers[serverid].MaxPlayers  = server.MaxPlayers
		case server.Mode != "":
			Data.Servers[serverid].Mode = server.Mode
		case server.PvP != "":
			Data.Servers[serverid].PvP = server.PvP
		case server.Version != "":
			Data.Servers[serverid].Version = server.Version
	}
} 