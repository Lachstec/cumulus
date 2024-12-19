package db

import (
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type ServerStore struct {
	db *sqlx.DB
}

func NewServerStore(db *sqlx.DB) Store[types.Server] {
	return &ServerStore{db: db}
}

func (s *ServerStore) GetById(id int64) (types.Server, error) {
	row := s.db.QueryRowx("SELECT * FROM mch_provisioner.servers WHERE id = $1;", id)
	var server types.Server
	err := row.StructScan(&server)

	if err != nil {
		return types.Server{}, err
	}

	return server, nil
}

func (s *ServerStore) Find(predicate Predicate[types.Server]) ([]types.Server, error) {
	rows, err := s.db.Queryx("SELECT * FROM mch_provisioner.servers;")
	if err != nil {
		return []types.Server{}, err
	}

	var servers []types.Server
	for rows.Next() {
		var server types.Server
		err = rows.StructScan(&server)
		if err != nil {
			return []types.Server{}, err
		}

		if predicate(server) {
			servers = append(servers, server)
		}
	}

	return servers, nil
}

func (s *ServerStore) Add(server types.Server) (int64, error) {
	var id int64
	err := s.db.QueryRowx(
		"INSERT INTO mch_provisioner.servers (openstack_id, name, addr, status, port, memory_mb, game, game_version, game_mode, difficulty, whitelist_enabled, players_max, ssh_key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id;",
		server.OpenstackId,
		server.Name,
		server.Address.String(),
		server.Status,
		server.Port,
		server.Memory,
		server.Game,
		server.GameVersion,
		server.GameMode,
		server.Difficulty,
		server.WhitelistEnabled,
		server.PlayersMax,
		server.SSHKey,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *ServerStore) Update(server types.Server) (types.Server, error) {
	var updated types.Server
	err := s.db.QueryRowx(
		"UPDATE mch_provisioner.servers SET openstack_id = $1, name = $2, addr = $3, status = $4, port = $5, memory_mb = $6, game = $7, game_version = $8, game_mode = $9, difficulty = $10, whitelist_enabled = $11, players_max = $12 WHERE id = $13 RETURNING *;",
		server.OpenstackId,
		server.Name,
		server.Address,
		server.Status,
		server.Port,
		server.Memory,
		server.Game,
		server.GameVersion,
		server.GameMode,
		server.Difficulty,
		server.WhitelistEnabled,
		server.PlayersMax,
		server.Id,
	).StructScan(&updated)
	if err != nil {
		return types.Server{}, err
	}
	return updated, nil
}

func (s *ServerStore) Delete(server types.Server) error {
	_, err := s.db.Exec("DELETE FROM mch_provisioner.servers WHERE id = $1", server.Id)
	if err != nil {
		return err
	}
	return nil
}
