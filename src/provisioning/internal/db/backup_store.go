package db

import (
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type ServerBackupStore struct {
	db sqlx.DB
}

func NewServerBackupStore(db sqlx.DB) *ServerBackupStore {
	return &ServerBackupStore{db: db}
}

func (b *ServerBackupStore) GetById(id int64) (types.Backup, error) {
	row := b.db.QueryRowx("SELECT * FROM mch_provisioner.world_backups WHERE id=?", id)
	var backup types.Backup
	err := row.StructScan(&backup)

	if err != nil {
		return types.Backup{}, err
	}
	return backup, nil
}

func (b *ServerBackupStore) Find(predicate Predicate[types.Backup]) ([]types.Backup, error) {
	rows, err := b.db.Queryx("SELECT * FROM mch_provisioner.world_backups")
	if err != nil {
		return []types.Backup{}, err
	}

	var backups []types.Backup
	for rows.Next() {
		var backup types.Backup
		err = rows.StructScan(&backup)

		if err != nil {
			return []types.Backup{}, err
		}

		if predicate(backup) {
			backups = append(backups, backup)
		}
	}

	return backups, nil
}

func (b *ServerBackupStore) Add(backup types.Backup) (types.Backup, error) {
	_, err := b.db.Exec("INSERT INTO mch_provisioner.world_backups (id, server_id, world, game, timestamp, size) VALUES (?, ?, ?, ?, ?, ?)",
		backup.Id,
		backup.ServerId,
		backup.World,
		backup.Game,
		backup.Timestamp,
		backup.Size,
		backup.Id,
	)
	if err != nil {
		return types.Backup{}, err
	}

	return backup, nil
}

func (b *ServerBackupStore) Update(backup types.Backup) (types.Backup, error) {
	_, err := b.db.Exec("UPDATE mch_provisioner.world_backups SET id = ?, server_id = ?, world = ?, game = ?, timestamp = ?, size = ? WHERE id = ?",
		backup.Id,
		backup.ServerId,
		backup.World,
		backup.Game,
		backup.Timestamp,
		backup.Size,
		backup.Id,
	)

	if err != nil {
		return types.Backup{}, err
	}

	return backup, nil
}

func (b *ServerBackupStore) Delete(backup types.Backup) error {
	_, err := b.db.Exec("DELETE FROM mch_provisioner.world_backups WHERE id=?", backup.Id)
	if err != nil {
		return err
	}
	return nil
}
