package services

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
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	ports2 "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/jmoiron/sqlx"
	"log"
	"strconv"
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
echo '{
	"mtu": 1442
}' > /etc/docker/daemon.json
systemctl restart docker
docker run -d -it -p 25565:25565 -v /mnt/data:/data -e EULA=TRUE itzg/minecraft-server
`

// MinecraftProvisioner is a service that can interact with
// OpenStack to provision and configure minecraft servers.
type MinecraftProvisioner struct {
	crypto      *CryptoService
	backupstore db.Store[types.Backup]
	serverstore db.Store[types.Server]
	keystore    db.Store[types.Key]
	ipstore     db.Store[types.FloatingIP]
	openstack   *openstack.Client
}

// NewMinecraftProvisioner creates a new provisioner service that stores
// its data in the database behind the given sqlx connection.
func NewMinecraftProvisioner(conn *sqlx.DB, openstack *openstack.Client, secretKey []byte) *MinecraftProvisioner {
	return &MinecraftProvisioner{
		crypto:      NewCryptoService(secretKey),
		backupstore: db.NewServerBackupStore(conn),
		serverstore: db.NewServerStore(conn),
		keystore:    db.NewKeyStore(conn),
		ipstore:     db.NewIPStore(conn),
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

// newKeyPair creates a new Public Key in Openstack that can be used
// to authenticate per SSH later down the line.
func (m *MinecraftProvisioner) newKeyPair(ctx context.Context, name string, publicKey string, privateKey string) (int64, error) {
	opts := keypairs.CreateOpts{
		Name:      name,
		PublicKey: publicKey,
	}

	client, err := m.openstack.ComputeClient()
	if err != nil {
		log.Println("Error getting compute client: ", err)
		return 0, err
	}

	keys, err := keypairs.Create(ctx, client, opts).Extract()
	if err != nil {
		log.Println("Error creating keypair: ", err)
		return 0, err
	}

	key := types.Key{
		Name:       keys.Name,
		PublicKey:  []byte(keys.PublicKey),
		PrivateKey: []byte(privateKey),
	}

	id, err := m.keystore.Add(&key)
	if err != nil {
		log.Println("Error adding keypair: ", err)
	}

	return id, nil
}

// makeFloatingIp creates a new Floating Ip for use to connect to the running game server.
// It currently has the network hardcoded to comply with the environment of the university cluster.
// Must be changed later down the line. The Ip gets automatically associated to the server given with serverId
func (m *MinecraftProvisioner) makeFloatingIp(ctx context.Context, serverId string) (*floatingips.FloatingIP, error) {
	client, err := m.openstack.NetworkingClient()
	if err != nil {
		fmt.Println("Error getting network client: ", err)
		return &floatingips.FloatingIP{}, err
	}

	compClient, err := m.openstack.ComputeClient()
	if err != nil {
		fmt.Println("Error getting compute client: ", err)
		return &floatingips.FloatingIP{}, err
	}

	// Wait for server to be in ACTIVE state
	for {
		server, err := servers.Get(ctx, compClient, serverId).Extract()
		if err != nil {
			log.Println("Error getting server status: ", err)
			return &floatingips.FloatingIP{}, err
		}

		if server.Status == "ACTIVE" {
			break
		}

		log.Println("Waiting for server to become ACTIVE...")
		time.Sleep(5 * time.Second) // Wait for 5 seconds before checking again
	}

	ports, err := ports2.List(client, ports2.ListOpts{
		DeviceID: serverId,
	}).AllPages(ctx)

	if err != nil {
		return &floatingips.FloatingIP{}, err
	}

	portList, err := ports2.ExtractPorts(ports)

	if len(portList) == 0 {
		fmt.Println("No ports found for device")
		return &floatingips.FloatingIP{}, err
	}

	fmt.Println("port id: ", portList[0].ID)

	ip, err := floatingips.Create(ctx, client, floatingips.CreateOpts{
		FloatingNetworkID: "6f530989-999a-49e6-9197-8a33ae7bfce7",
		PortID:            portList[0].ID,
	}).Extract()

	if err != nil {
		fmt.Println("Error creating floating ip: ", err)
		return &floatingips.FloatingIP{}, err
	}

	return ip, nil
}

// WaitForVolumeReady wait until a requested volume is ready to be mounted
// Returns when the volume is ready or errors if the given timeout has elapsed.
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
func (m *MinecraftProvisioner) NewGameServer(ctx context.Context, server *types.Server, user *types.User) (*types.Server, error) {
	client, err := m.openstack.ComputeClient()
	if err != nil {
		log.Println("Error getting compute client: ", err)
		return nil, err
	}

	volume, err := m.newPersistentVolume(ctx, server.Name+"_volume")
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
			UUID:                string(server.Image),
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

	kid, err := m.newKeyPair(ctx, server.Name+"public_key", publicKey, privateKey)
	if err != nil {
		log.Println("Error saving pubkey to openstack: ", err)
		return nil, err
	}

	opts := servers.CreateOpts{
		Name:        server.Name,
		FlavorRef:   strconv.FormatInt(types.Flavours[server.Flavour-1].ID, 10),
		ImageRef:    string(server.Image),
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
		KeyName:           server.Name + "public_key",
	}

	gcServer, err := servers.Create(ctx, client, optsExt, nil).Extract()
	if err != nil {
		log.Println("Error spawning server: ", err)
		return nil, err
	}

	addr, err := m.makeFloatingIp(ctx, gcServer.ID)
	if err != nil {
		log.Println("Error creating floating ip: ", err)
		return nil, err
	}

	ip, err := m.ipstore.Add(&types.FloatingIP{
		OpenstackId: addr.ID,
		Ip:          addr.FloatingIP,
	})

	if err != nil {
		log.Println("Error adding floating ip: ", err)
		return nil, err
	}

	server.OpenstackId = gcServer.ID
	server.Address = ip
	server.Status = types.Running
	server.Port = 25565
	server.SSHKey = kid
	server.OpenstackId = gcServer.ID
	server.UserID = user.ID

	id, err := m.serverstore.Add(server)
	if err != nil {
		log.Println("Error adding server: ", err)
		return nil, err
	}

	server.ID = id

	backup := &types.Backup{
		OpenstackID: volume,
		ServerID:    server.ID,
		Timestamp:   time.Now(),
		Size:        10000,
	}

	_, err = m.backupstore.Add(backup)

	if err != nil {
		return nil, err
	}

	return server, nil
}

// DeleteGameServer completely de-provisions the given server. It does this by
// first deleting the compute instance, then the keypair associated with it. After that,
// the attached volume with game data gets deleted and as a last thing, the floating
// ip that was used to make it accessible gets released.
func (m *MinecraftProvisioner) DeleteGameServer(ctx context.Context, server types.Server) error {
	backups, err := m.backupstore.Find(func(b *types.Backup) bool { return b.ServerID == server.ID })
	if err != nil {
		log.Printf("Error finding backups for Server with id %d: %v", server.ID, err)
		return err
	}

	storageClient, err := m.openstack.StorageClient()
	if err != nil {
		log.Println("Error getting storage client: ", err)
		return err
	}

	computeClient, err := m.openstack.ComputeClient()
	if err != nil {
		log.Println("Error getting compute client: ", err)
		return err
	}

	networkClient, err := m.openstack.NetworkingClient()
	if err != nil {
		log.Println("Error getting network client: ", err)
		return err
	}

	key, err := m.keystore.Find(func(k *types.Key) bool { return server.SSHKey == k.Id })
	if err != nil || len(key) == 0 {
		log.Println("Error finding server key: ", err)
	}

	err = servers.Delete(ctx, computeClient, server.OpenstackId).ExtractErr()
	if err != nil {
		log.Println("Error deleting server: ", err)
		return err
	}

	err = keypairs.Delete(ctx, computeClient, key[0].Name, nil).ExtractErr()
	if err != nil {
		log.Println("Error deleting keypair: ", err)
	}

	for _, backup := range backups {
		log.Println("Deleting backup: ", backup.OpenstackID)
		err = volumes.Delete(ctx, storageClient, backup.OpenstackID, nil).ExtractErr()
		if err != nil {
			log.Println("Error deleting volume: ", err)
		}
		log.Println("Deleting backup: ", backup.OpenstackID)
		err = m.WaitForVolumeReady(ctx, backup.OpenstackID, time.Minute*2)
		if err != nil {
			log.Println("Error deleting backup: ", err)
			return err
		}

		err = volumes.Delete(ctx, storageClient, backup.OpenstackID, nil).ExtractErr()
		if err != nil {
			log.Println("Error deleting backup: ", err)
			return err
		}
	}

	ip, err := m.ipstore.Find(func(ip *types.FloatingIP) bool { return server.Address == ip.Id })
	if err != nil || len(ip) == 0 {
		log.Println("Error finding ip: ", err)
		return err
	}

	err = floatingips.Delete(ctx, networkClient, ip[0].OpenstackId).ExtractErr()
	if err != nil {
		log.Println("Error deleting floating ip: ", err)
		return err
	}

	return nil
}
