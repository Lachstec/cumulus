package services

import (
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type ServerService struct {
	store db.Store[types.Server]
}

func NewServerService(conn *sqlx.DB) *ServerService {
	return &ServerService{
		store: db.NewServerStore(conn),
	}
}

func (c *ServerService) ReadAllServers() ([]*types.Server, error) {

	servers, err := c.store.Find(func(s *types.Server) bool { return true })
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (c *ServerService) ReadServerByServerID(serverid int64) ([]*types.Server, error) {
	server, err := c.store.Find(func(s *types.Server) bool { return s.ID == serverid })
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (c *ServerService) CreateServer(server *types.Server) (int64, error) {
	serverid, err := c.store.Add(server)
	if err != nil {
		return 0, err
	}
	return serverid, nil
}

func (c *ServerService) DeleteServer(server *types.Server) (error) {
	err := c.store.Delete(server)
	if err != nil {
		return err
	}
	return nil
}

func (c *ServerService) UpdateServer(server *types.Server) (*types.Server, error) {
	server, err := c.store.Update(server)
	if err != nil {
		return nil, err
	}
	return server, nil
}
