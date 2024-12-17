package openstack

import (
	"context"
	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack"
)

// Client is a factory struct that can construct Clients
// for the various OpenStack Services. The Clients returned from it all
// use the configuration that was used to create the Client by calling NewClient.
type Client struct {
	options gophercloud.AuthOptions
	client  *gophercloud.ProviderClient
}

// NewClient creates a new Client that uses the
// account specified in config to connect to OpenStack.
func NewClient(config config.Config) (*Client, error) {
	authOpts := config.Openstack.AuthOptions()
	ctx := context.Background()
	provider, err := openstack.AuthenticatedClient(ctx, authOpts)

	if err != nil {
		return nil, err
	}

	client := &Client{
		options: config.Openstack.AuthOptions(),
		client:  provider,
	}

	return client, nil
}

func (c *Client) ComputeClient() (*gophercloud.ServiceClient, error) {
	client, err := openstack.NewComputeV2(c.client, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) StorageClient() (*gophercloud.ServiceClient, error) {
	client, err := openstack.NewBlockStorageV3(c.client, gophercloud.EndpointOpts{})
	if err != nil {
		return nil, err
	}
	return client, nil
}
