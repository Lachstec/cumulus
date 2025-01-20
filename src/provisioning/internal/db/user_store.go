package db

import (
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sqlx.DB) Store[types.User] {
	return &UserStore{db: db}
}

func (s *UserStore) GetById(ID int64) (*types.User, error) {
	row := s.db.QueryRowx("SELECT * FROM mch_provisioner.users WHERE ID = $1;", ID)
	var user *types.User
	err := row.StructScan(&user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStore) Find(predicate Predicate[*types.User]) ([]*types.User, error) {
	rows, err := s.db.Queryx("SELECT * FROM mch_provisioner.users;")
	if err != nil {
		return nil, err
	}

	var users []*types.User
	for rows.Next() {
		var user *types.User
		err = rows.StructScan(&user)
		if err != nil {
			return nil, err
		}

		if predicate(user) {
			users = append(users, user)
		}
	}

	return users, nil
}

func (s *UserStore) Add(user *types.User) (int64, error) {
	var ID int64
	err := s.db.QueryRowx(
		"INSERT INTO mch_provisioner.users (sub, name, class) VALUES ($1, $2, $3) RETURNING ID;",
		user.Sub,
		user.Name,
		user.Class,
	).Scan(&ID)
	if err != nil {
		return 0, err
	}
	return ID, nil
}

func (s *UserStore) Update(user *types.User) (*types.User, error) {
	var updated *types.User
	err := s.db.QueryRowx(
		"UPDATE mch_provisioner.users SET sub = $1, name = $2, class = $3, WHERE ID = $4 RETURNING *;",
		user.Sub,
		user.Name,
		user.Class,
		user.ID,
	).StructScan(&updated)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *UserStore) Delete(user *types.User) error {
	_, err := s.db.Exec("DELETE FROM mch_provisioner.users WHERE ID = $1", user.ID)
	if err != nil {
		return err
	}
	return nil
}
