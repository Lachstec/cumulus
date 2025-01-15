package types

import "net"

type FloatingIP struct {
	Id          int64
	OpenstackId string `db:"openstack_id"`
	Ip          string `db:"addr"`
}

func (f *FloatingIP) GetIP() net.IP {
	return net.ParseIP(f.Ip)
}
