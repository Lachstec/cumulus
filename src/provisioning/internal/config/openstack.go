package config

import (
	"github.com/gophercloud/gophercloud/v2"
)

// OpenStackConfig contains all information that is needed to authenticate
// to an openstack cluster.
type OpenStackConfig struct {
	// identityEndpoint URL to the Keystone Service of the openstack cluster
	identityEndpoint string
	// username to use for authentication
	username string
	// password to use for authentication
	password string
	// tenantId to use for authentication
	tenantId string
}

// AuthOptions return an instance of gophercloud.AuthOptions that can be used
// to authenticate against openstacks keystone service.
func (c OpenStackConfig) AuthOptions() gophercloud.AuthOptions {
	return gophercloud.AuthOptions{
		IdentityEndpoint: c.identityEndpoint,
		Username:         c.username,
		Password:         c.password,
		TenantID:         c.tenantId,
	}
}
