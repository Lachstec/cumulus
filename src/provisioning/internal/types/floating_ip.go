package types

import "net"

type FloatingIP struct {
	Id          int64
	OpenstackId string `db:"openstack_id"`
	Ip          net.IP `db:"addr"`
}
