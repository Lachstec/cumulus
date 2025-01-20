package db

import (
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type KeyStore struct {
	db *sqlx.DB
}

func NewKeyStore(db *sqlx.DB) Store[types.Key] {
	return &KeyStore{db: db}
}

func (k *KeyStore) GetById(id int64) (*types.Key, error) {
	row := k.db.QueryRowx("SELECT * FROM mch_provisioner.keypairs WHERE id=$1;", id)
	var key types.Key
	err := row.StructScan(&key)

	if err != nil {
		return &types.Key{}, err
	}
	return &key, nil
}

func (k *KeyStore) Find(predicate Predicate[*types.Key]) ([]*types.Key, error) {
	rows, err := k.db.Queryx("SELECT * FROM mch_provisioner.keypairs;")
	if err != nil {
		return []*types.Key{}, err
	}

	var keys []*types.Key
	for rows.Next() {
		var key types.Key
		err = rows.StructScan(&key)

		if err != nil {
			return []*types.Key{}, err
		}

		if predicate(&key) {
			keys = append(keys, &key)
		}
	}

	return keys, nil
}

func (k *KeyStore) Add(key *types.Key) (int64, error) {
	var id int64

	err := k.db.QueryRowx("INSERT INTO mch_provisioner.keypairs (name, public_key, private_key) VALUES ($1, $2, $3) RETURNING id;",
		key.Name, key.PublicKey, key.PrivateKey).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (k *KeyStore) Update(key *types.Key) (*types.Key, error) {
	_, err := k.db.Exec("UPDATE mch_provisioner.keypairs SET name = $1, public_key = $2, private_key = $3 WHERE id = $4;",
		key.Name, key.PublicKey, key.PrivateKey, key.Id)
	if err != nil {
		return &types.Key{}, err
	}

	return key, nil
}

func (k *KeyStore) Delete(key *types.Key) error {
	_, err := k.db.Exec("DELETE FROM mch_provisioner.keypairs WHERE id=$1;", key.Id)
	if err != nil {
		return err
	}
	return nil
}
