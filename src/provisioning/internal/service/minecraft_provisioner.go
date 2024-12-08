package service

import (
	"context"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/openstack"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/jmoiron/sqlx"
	"net"
)

// MinecraftProvisioner is a service that can interact with
// OpenStack to provision and configure minecraft servers.
type MinecraftProvisioner struct {
	store     db.Store[types.Server]
	openstack openstack.Client
}

// NewMinecraftProvisioner creates a new provisioner service that stores
// its data in the database behind the given sqlx connection.
func NewMinecraftProvisioner(conn *sqlx.DB, openstack openstack.Client) *MinecraftProvisioner {
	return &MinecraftProvisioner{
		store:     db.NewServerStore(conn),
		openstack: openstack,
	}
}

// NewGameServer provisions a new Gameserver with the specified flavour in openstack. The provisioned server
// has an ephemeral disk and uses the default settings and config of the specified image
// in openstack. Information about the server gets stored in the database.
func (m *MinecraftProvisioner) NewGameServer(ctx context.Context, name string, flavour types.Flavour, image types.Image) (*types.Server, error) {
	client, err := m.openstack.ComputeClient()
	if err != nil {
		return nil, err
	}

	blockDev := []servers.BlockDevice{
		{
			BootIndex:           0,
			DeleteOnTermination: true,
			DestinationType:     servers.DestinationLocal,
			SourceType:          servers.SourceImage,
			UUID:                image.Value(),
		},
	}

	opts := servers.CreateOpts{
		Name:        name,
		FlavorRef:   flavour.Value(),
		ImageRef:    image.Value(),
		BlockDevice: blockDev,
	}

	server, err := servers.Create(ctx, client, opts, nil).Extract()
	if err != nil {
		return nil, err
	}

	gameserver := types.Server{
		Name:             name,
		Address:          net.ParseIP(server.AccessIPv4),
		Status:           types.Running,
		Port:             25565,
		Memory:           flavour.AvailableRam(),
		Game:             "Minecraft",
		GameVersion:      "1.0.0",
		GameMode:         types.Survival,
		Difficulty:       types.Normal,
		WhitelistEnabled: false,
		PlayersMax:       2,
	}

	id, err := m.store.Add(gameserver)
	if err != nil {
		return nil, err
	}

	gameserver, err = m.store.GetById(id)
	if err != nil {
		return nil, err
	}

	return &gameserver, nil
}
