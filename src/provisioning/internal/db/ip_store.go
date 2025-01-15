package db

import (
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
)

type IPStore struct {
	db *sqlx.DB
}

func NewIPStore(db *sqlx.DB) Store[types.FloatingIP] {
	return &IPStore{db: db}
}

func (i *IPStore) GetById(id int64) (types.FloatingIP, error) {
	row := i.db.QueryRowx("SELECT * FROM mch_provisioner.floating_ips WHERE id = $1;", id)
	var ip types.FloatingIP

	err := row.StructScan(&ip)
	if err != nil {
		return types.FloatingIP{}, err
	}

	return ip, nil
}

func (i *IPStore) Find(predicate Predicate[types.FloatingIP]) ([]types.FloatingIP, error) {
	rows, err := i.db.Queryx("SELECT * FROM mch_provisioner.floating_ips;")
	if err != nil {
		return []types.FloatingIP{}, err
	}

	var ips []types.FloatingIP
	for rows.Next() {
		var ip types.FloatingIP
		err = rows.StructScan(&ip)

		if err != nil {
			return []types.FloatingIP{}, err
		}

		if predicate(ip) {
			ips = append(ips, ip)
		}
	}

	return ips, nil
}

func (i *IPStore) Add(ip types.FloatingIP) (int64, error) {
	var id int64

	err := i.db.QueryRowx("INSERT INTO mch_provisioner.floating_ips (openstack_id, addr) VALUES ($1, $2) RETURNING id;",
		ip.OpenstackId, ip.Ip).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (i *IPStore) Update(ip types.FloatingIP) (types.FloatingIP, error) {
	_, err := i.db.Exec("UPDATE mch_provisioner.floating_ips SET openstack_id = $1, addr = $2 WHERE id = $3;",
		ip.OpenstackId, ip.Ip, ip.Id)

	if err != nil {
		return types.FloatingIP{}, err
	}

	return ip, nil
}

func (i *IPStore) Delete(ip types.FloatingIP) error {
	_, err := i.db.Exec("DELETE FROM mch_provisioner.floating_ips WHERE id = $1;", ip.Id)
	if err != nil {
		return err
	}
	return nil
}
