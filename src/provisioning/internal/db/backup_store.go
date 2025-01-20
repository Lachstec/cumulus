package db

import (
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type ServerBackupStore struct {
	db *sqlx.DB
}

func NewServerBackupStore(db *sqlx.DB) Store[types.Backup] {
	return &ServerBackupStore{db: db}
}

func (b *ServerBackupStore) GetById(ID int64) (*types.Backup, error) {
	row := b.db.QueryRowx("SELECT * FROM mch_provisioner.world_backups WHERE ID=$1", ID)
	var backup *types.Backup
	err := row.StructScan(&backup)

	if err != nil {
		return nil, err
	}
	return backup, nil
}

func (b *ServerBackupStore) Find(predicate Predicate[*types.Backup]) ([]*types.Backup, error) {
	rows, err := b.db.Queryx("SELECT * FROM mch_provisioner.world_backups")
	if err != nil {
		return nil, err
	}

	var backups []*types.Backup
	for rows.Next() {
		var backup *types.Backup
		err = rows.StructScan(&backup)

		if err != nil {
			return nil, err
		}

		if predicate(backup) {
			backups = append(backups, backup)
		}
	}

	return backups, nil
}

func (b *ServerBackupStore) Add(backup *types.Backup) (int64, error) {
	var ID int64
	err := b.db.QueryRowx("INSERT INTO mch_provisioner.world_backups (openstack_id, server_id, timestamp, size) VALUES ($1, $2, $3, $4) RETURNING id;",
		backup.OpenstackID,
		backup.ServerID,
		backup.Timestamp,
		backup.Size,
	).Scan(&ID)
	if err != nil {
		return 0, err
	}

	return ID, nil
}

func (b *ServerBackupStore) Update(backup *types.Backup) (*types.Backup, error) {
	_, err := b.db.Exec("UPDATE mch_provisioner.world_backups SET server_id = $1, timestamp = $2, size = $3 WHERE id = $4",
		backup.ServerID,
		backup.Timestamp,
		backup.Size,
		backup.ID,
	)

	if err != nil {
		return nil, err
	}

	return backup, nil
}

func (b *ServerBackupStore) Delete(backup *types.Backup) error {
	_, err := b.db.Exec("DELETE FROM mch_provisioner.world_backups WHERE ID = $1", backup.ID)
	if err != nil {
		return err
	}
	return nil
}
