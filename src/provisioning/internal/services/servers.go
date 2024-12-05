package services

import (
	"github.com/Lachstec/mc-hosting/internal/types"
)

func ReadNumOfServers() int {
	return len(types.Servers)
}

func ReadAllServers() []types.Server {
	return types.Servers
}

func ReadServerByServerID(serverid int) types.Server {
	return types.Servers[serverid]
}

func CreateServer(server types.Server) {
	types.Servers = append(types.Servers, server)
}

func DeleteServerByServerID(serverid int) {
	types.Servers = append(types.Servers[:serverid], types.Servers[serverid+1:]...)
}

func UpdateServer(serverid int, server types.Server) {
	switch {
	case server.Name != "":
		types.Servers[serverid].Name = server.Name
	case server.Difficulty != "":
		types.Servers[serverid].Difficulty = server.Difficulty
	case server.IP != "":
		types.Servers[serverid].IP = server.IP
	case server.MaxPlayers != 0:
		types.Servers[serverid].MaxPlayers = server.MaxPlayers
	case server.Mode != "":
		types.Servers[serverid].Mode = server.Mode
	case server.PvP != "":
		types.Servers[serverid].PvP = server.PvP
	case server.Version != "":
		types.Servers[serverid].Version = server.Version
	}
}
