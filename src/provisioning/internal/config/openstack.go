package config

import (
	"github.com/gophercloud/gophercloud/v2"
)

// OpenStackConfig contains all information that is needed to authenticate
// to an openstack cluster.
type OpenStackConfig struct {
	// IdentityEndpoint URL to the Keystone Service of the openstack cluster
	IdentityEndpoint string
	// Username to use for authentication
	Username string
	// Password to use for authentication
	Password string
	// Domain to use for authentication
	Domain string
	// TenantName to use for authentication
	TenantName string
}

// AuthOptions return an instance of gophercloud.AuthOptions that can be used
// to authenticate against openstacks keystone service.
func (c OpenStackConfig) AuthOptions() gophercloud.AuthOptions {
	return gophercloud.AuthOptions{
		IdentityEndpoint: c.IdentityEndpoint,
		Username:         c.Username,
		Password:         c.Password,
		DomainName:       c.Domain,
		TenantName:       c.TenantName,
	}
}
