package db

import (
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type ServerStore struct {
	db *sqlx.DB
}

func NewServerStore(db *sqlx.DB) *ServerStore {
	return &ServerStore{db: db}
}

func (s *ServerStore) GetById(id int64) (types.Server, error) {
	row := s.db.QueryRowx("SELECT * FROM mch_provisioner.servers WHERE id = ?;", id)
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

func (s *ServerStore) Add(server types.Server) (types.Server, error) {
	_, err := s.db.Exec(
		"INSERT INTO mch_provisioner.servers (id, addr, status, port, memory_mb, game, game_version, game_mode, difficulty, whitelist_enabled, players_max) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);",
		server.Id,
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
	)
	if err != nil {
		return types.Server{}, err
	}
	return server, nil
}

func (s *ServerStore) Update(server types.Server) (types.Server, error) {
	_, err := s.db.Exec(
		"UPDATE mch_provisioner.servers SET id = ?, addr = ?, status = ?, port = ?, memory_mb = ?, game = ?, game_version = ?, game_mode = ?, difficulty = ?, whitelist_enabled = ?, players_max = ? WHERE id = ?;",
		server.Id,
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
	)
	if err != nil {
		return types.Server{}, err
	}
	return server, nil
}

func (s *ServerStore) Delete(server types.Server) error {
	_, err := s.db.Exec("DELETE FROM mch_provisioner.servers WHERE id = ?", server.Id)
	if err != nil {
		return err
	}
	return nil
}
