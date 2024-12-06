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

func (c *ServerService) UpdateServer(serverid int64, server types.Server) (types.Server, error) {
	server.ID = serverid
	server, err := c.store.Update(server)
	if err != nil {
		return types.Nothing[types.Server](), err
	}
	return server, nil
}
