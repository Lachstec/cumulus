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

func (s *ServerStore) GetById(ID int64) (*types.Server, error) {
	row := s.db.QueryRowx("SELECT * FROM mch_provisioner.servers WHERE ID = $1;", ID)
	var server types.Server
	err := row.StructScan(&server)

	if err != nil {
		return nil, err
	}

	return &server, nil
}

func (s *ServerStore) Find(predicate Predicate[*types.Server]) ([]*types.Server, error) {
	rows, err := s.db.Queryx("SELECT * FROM mch_provisioner.servers;")
	if err != nil {
		return []*types.Server{}, err
	}

	var servers []*types.Server
	for rows.Next() {
		var server types.Server
		err = rows.StructScan(&server)
		if err != nil {
			return nil, err
		}

		if predicate(&server) {
			servers = append(servers, &server)
		}
	}
	return servers, nil
}

func (s *ServerStore) Add(server *types.Server) (int64, error) {
	var ID int64
	err := s.db.QueryRowx(
		"INSERT INTO mch_provisioner.servers(userid, openstack_id, name, addr, status, port, flavour, image, game, game_version, game_mode, difficulty, whitelist_enabled, pvp_enabled, players_max, ssh_key) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id;",
		server.UserID,
		server.OpenstackId,
		server.Name,
		server.Address,
		server.Status,
		server.Port,
		server.Flavour,
		server.Image,
		server.Game,
		server.GameVersion,
		server.GameMode,
		server.Difficulty,
		server.WhitelistEnabled,
		server.PvPEnabled,
		server.PlayersMax,
		server.SSHKey,
	).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (s *ServerStore) Update(server *types.Server) (*types.Server, error) {
	var updated types.Server
	err := s.db.QueryRowx(
		"UPDATE mch_provisioner.servers SET openstack_id = $1, name = $2, addr = $3, status = $4, port = $5, flavour = $6, image = $7, game = $8, game_version = $9, game_mode = $10, difficulty = $11, whitelist_enabled = $12, pvp_enabled = $13, players_max = $14 WHERE id = $15 RETURNING *;",
		server.OpenstackId,
		server.Name,
		server.Address,
		server.Status,
		server.Port,
		server.Flavour,
		server.Image,
		server.Game,
		server.GameVersion,
		server.GameMode,
		server.Difficulty,
		server.WhitelistEnabled,
		server.PvPEnabled,
		server.PlayersMax,
		server.ID,
	).StructScan(&updated)
	if err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *ServerStore) Delete(server *types.Server) error {
	_, err := s.db.Exec("DELETE FROM mch_provisioner.servers WHERE ID = $1", server.ID)
	if err != nil {
		return err
	}
	return nil
}
