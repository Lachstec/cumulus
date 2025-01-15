package types

import "net"

// FloatingIP represents a Floating IP that can be associated to a types.Server.
type FloatingIP struct {
	// Id is the unique identifier for the FloatingIP in the database.
	Id int64
	// OpenstackId is the ID of the FloatingIP in Open Stack
	OpenstackId string `db:"openstack_id"`
	// Ip contains the actual outside-facing IP Address
	Ip string `db:"addr"`
}

// GetIP returns a net.IP containing the stored IP Address.
func (f *FloatingIP) GetIP() net.IP {
	return net.ParseIP(f.Ip)
}
