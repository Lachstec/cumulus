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
		"INSERT INTO mch_provisioner.servers (name, addr, status, port, memory_mb, game, game_version, game_mode, difficulty, whitelist_enabled, players_max, ssh_key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;",
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
		"UPDATE mch_provisioner.servers SET name = $1, addr = $2, status = $3, port = $4, memory_mb = $5, game = $6, game_version = $7, game_mode = $8, difficulty = $9, whitelist_enabled = $10, players_max = $11 WHERE id = $12 RETURNING *;",
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
