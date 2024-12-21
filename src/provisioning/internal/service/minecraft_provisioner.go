package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/openstack"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servers"
	"github.com/jmoiron/sqlx"
	"log"
	"net"
	"time"
)

// initScript is a simple and more or less dangerous way of getting docker ready to roll
// on newly provisioned servers. This should NOT be used in a production scenario
// as we are piping a more or less random script into bash. Installation should
// be performed with the package manager of the used distribution when deploying to production.
const initScript = `#!/bin/bash
device="/dev/vdb"
mount_point="/mnt/data"

curl -fsSL get.docker.com -o get-docker.sh && sh get-docker.sh

for i in {1..24}; do
	if lsblk | grep -q "$(basename $device)"; then
		echo "Device available"
		break
	fi
	echo "Waiting for Device to be available"
	sleep 5
done

if [ ! -b "$device" ]; then
	echo "Error: Device not available"
fi

if [ -z "$(blkid $device)" ]; then
    echo "Formatting $device as ext4..."
    mkfs.ext4 $device
fi

mkdir -p $mount_point
mount $device $mount_point
echo "$device $mount_point ext4 defaults 0 0" >> /etc/fstab

docker run -d -it -p 25565:25565 -e EULA=TRUE itzg/minecraft-server
`

// MinecraftProvisioner is a service that can interact with
// OpenStack to provision and configure minecraft servers.
type MinecraftProvisioner struct {
	crypto      *CryptoService
	backupstore db.Store[types.Backup]
	serverstore db.Store[types.Server]
	openstack   *openstack.Client
}

// NewMinecraftProvisioner creates a new provisioner service that stores
// its data in the database behind the given sqlx connection.
func NewMinecraftProvisioner(conn *sqlx.DB, openstack *openstack.Client, secretKey []byte) *MinecraftProvisioner {
	return &MinecraftProvisioner{
		crypto:      NewCryptoService(secretKey),
		backupstore: db.NewServerBackupStore(conn),
		serverstore: db.NewServerStore(conn),
		openstack:   openstack,
	}
}

// newPersistentVolume creates a new, persisting volume to store minecraft world data in.
// Returns the ID of the newly created volume or an error.
func (m *MinecraftProvisioner) newPersistentVolume(ctx context.Context, name string) (string, error) {
	opts := volumes.CreateOpts{
		Name: name,
		Size: 10,
	}

	client, err := m.openstack.StorageClient()
	if err != nil {
		log.Println("Error getting storage client: ", err)
		return "", err
	}

	vol, err := volumes.Create(ctx, client, opts, nil).Extract()
	if err != nil {
		log.Println("Error creating volume: ", err)
		return "", err
	}

	return vol.ID, nil
}

func (m *MinecraftProvisioner) newKeyPair(ctx context.Context, name string, publicKey string) error {
	opts := keypairs.CreateOpts{
		Name:      name,
		PublicKey: publicKey,
	}

	client, err := m.openstack.ComputeClient()
	if err != nil {
		log.Println("Error getting compute client: ", err)
		return err
	}

	_, err = keypairs.Create(ctx, client, opts).Extract()
	if err != nil {
		log.Println("Error creating keypair: ", err)
		return err
	}
	return nil
}

func (m *MinecraftProvisioner) pollServerIp(ctx context.Context, serverId string) (string, error) {
	timeout := time.After(2 * time.Minute)
	ticker := time.Tick(5 * time.Second)
	client, err := m.openstack.ComputeClient()
	if err != nil {
		log.Println("Error getting compute client: ", err)
		return "", err
	}

	for {
		select {
		case <-timeout:
			return "", fmt.Errorf("timed out waiting for server IP")
		case <-ticker:
			server, err := servers.Get(ctx, client, serverId).Extract()
			if err != nil {
				return "", fmt.Errorf("failed to get server details: %v", err)
			}

			// Check for IP address
			for _, addresses := range server.Addresses {
				for _, addr := range addresses.([]interface{}) {
					address := addr.(map[string]interface{})
					if address["addr"] != nil {
						return address["addr"].(string), nil
					}
				}
			}
		}
	}
}

func (m *MinecraftProvisioner) WaitForVolumeReady(ctx context.Context, volumeID string, timeout time.Duration) error {
	start := time.Now()
	client, err := m.openstack.StorageClient()
	if err != nil {
		log.Println("Error getting compute client: ", err)
		return err
	}
	for time.Since(start) < timeout {
		volume, err := volumes.Get(ctx, client, volumeID).Extract()
		if err != nil {
			fmt.Printf("Volume not ready yet, retrying: %v\n", err)
			time.Sleep(5 * time.Second) // Wait before retrying
			continue
		}

		// Check if volume is ready
		if volume.Status == "available" {
			fmt.Printf("Volume %s is now ready.\n", volumeID)
			return nil
		}

		fmt.Printf("Volume status: %s, waiting...\n", volume.Status)
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("volume %s did not become ready within timeout", volumeID)
}

// NewGameServer provisions a new Gameserver with the specified flavour in openstack. The provisioned server
// has an ephemeral disk and uses the default settings and config of the specified image
// in openstack. Information about the server gets stored in the database.
func (m *MinecraftProvisioner) NewGameServer(ctx context.Context, name string, flavour types.Flavour, image types.Image) (*types.Server, error) {
	client, err := m.openstack.ComputeClient()
	if err != nil {
		log.Println("Error getting compute client: ", err)
		return nil, err
	}

	volume, err := m.newPersistentVolume(ctx, name+"_volume")
	if err != nil {
		log.Println("Error creating persistent volume: ", err)
		return nil, err
	}

	err = m.WaitForVolumeReady(ctx, volume, 2*time.Minute)
	if err != nil {
		log.Println("Error waiting for volume to become available: ", err)
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
		{
			BootIndex:           1,
			DeleteOnTermination: false,
			SourceType:          servers.SourceVolume,
			DestinationType:     servers.DestinationVolume,
			UUID:                volume,
		},
	}

	userData := base64.StdEncoding.EncodeToString([]byte(initScript))

	publicKey, privateKey, err := m.crypto.NewKeyPair()
	if err != nil {
		log.Println("Error generating keys for server: ", err)
		return nil, err
	}

	err = m.newKeyPair(ctx, name+"public_key", publicKey)
	if err != nil {
		log.Println("Error saving pubkey to openstack: ", err)
		return nil, err
	}

	opts := servers.CreateOpts{
		Name:        name,
		FlavorRef:   flavour.Value(),
		ImageRef:    image.Value(),
		BlockDevice: blockDev,
		UserData:    []byte(userData),
		Networks: []servers.Network{
			{
				UUID: "9efbb5f1-ff47-45f4-9d06-77873aff7eb4",
			},
		},
	}

	optsExt := keypairs.CreateOptsExt{
		CreateOptsBuilder: opts,
		KeyName:           name + "public_key",
	}

	server, err := servers.Create(ctx, client, optsExt, nil).Extract()
	if err != nil {
		log.Println("Error spawning server: ", err)
		return nil, err
	}

	addr, err := m.pollServerIp(ctx, server.ID)
	if err != nil {
		log.Println("Error polling server IP address: ", err)
	}

	gameserver := types.Server{
		OpenstackId:      server.ID,
		Name:             name,
		Address:          net.ParseIP(addr),
		Status:           types.Running,
		Port:             25565,
		Memory:           flavour.AvailableRam(),
		Game:             "Minecraft",
		GameVersion:      "1.0.0",
		GameMode:         types.Survival,
		Difficulty:       types.Normal,
		WhitelistEnabled: false,
		PlayersMax:       2,
		SSHKey:           []byte(privateKey),
	}

	id, err := m.serverstore.Add(gameserver)
	if err != nil {
		log.Println("Error adding server to database: ", err)
		return nil, err
	}

	gameserver, err = m.serverstore.GetById(id)
	if err != nil {
		log.Println("Error getting server from database: ", err)
		return nil, err
	}

	backup := types.Backup{
		OpenstackId: volume,
		ServerId:    gameserver.Id,
		Timestamp:   time.Now(),
		Size:        10000,
	}

	_, err = m.backupstore.Add(backup)

	if err != nil {
		log.Println("Error adding backup to database: ", err)
	}

	return &gameserver, nil
}
