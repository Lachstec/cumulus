package services

import (
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type ServerService struct {
	store *db.ServerStore
}

func NewServerService(conn *sqlx.DB) *ServerService {
	return &ServerService{
		store: db.NewServerStore(conn),
	}
}

func ReadNumOfServers() int {
	return len(types.Servers)
}

func (c *ServerService) ReadAllServers() ([]types.Server, error) {

	servers, err := c.store.Find(func(s types.Server) bool { return true })
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (c *ServerService) ReadServerByServerID(serverid int64) ([]types.Server, error) {
	server, err := c.store.Find(func(s types.Server) bool { return s.ID == serverid })
	if err != nil {
		return nil, err
	}
	return server, nil
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
