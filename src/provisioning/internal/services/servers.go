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

func (c *ServerService) ReadAllServers() ([]types.Server, error) {

	servers, err := c.store.Find(func(s types.Server) bool { return true })
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (c *ServerService) ReadServerByServerID(serverid int64) (types.Server, error) {
	server, err := c.store.Find(func(s types.Server) bool { return s.ID == serverid })
	if err != nil {
		return types.Nothing[types.Server](), err
	}
	return server[0], nil
}

func (c *ServerService) CreateServer(server types.Server) (int64, error) {
	serverid, err := c.store.Add(server)
	if err != nil {
		return 0, err
	}
	return serverid, nil
}

func (c *ServerService) DeleteServerByServerID(serverid int64) (error) {
	server, err := c.ReadServerByServerID(serverid)
	if err != nil {
		return err
	}
	err = c.store.Delete(server)
	if err != nil {
		return err
	}
	return nil
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
