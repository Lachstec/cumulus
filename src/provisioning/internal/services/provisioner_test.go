//go:build exclude

package service

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"github.com/Lachstec/mc-hosting/internal/config"
	"github.com/Lachstec/mc-hosting/internal/db"
	openstack2 "github.com/Lachstec/mc-hosting/internal/openstack"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/ssh"
	"math/rand"
	"testing"
)

type TestConnection struct {
	Db *sqlx.DB
}

func NewTestConnection() *TestConnection {
	conn, err := sqlx.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic("unable to connect to testing database")
	}
	migrator := db.NewMigrator(conn)
	err = migrator.Migrate("../../migrations")
	if err != nil {
		panic(err)
	}

	return &TestConnection{conn}
}

func TestServerCreation(t *testing.T) {
	// Just a random key, don't use for production!
	key, err := base64.StdEncoding.DecodeString("1YRCJE3rUygZv4zXUhBNUf1sDUIszdT2KAtczVYB85c=")
	if err != nil {
		t.Fatal(err)
	}
	cfg := config.Config{
		Db:    config.DbConfig{},
		Auth0: config.Auth0Config{},
		Openstack: config.OpenStackConfig{
			IdentityEndpoint: "<keystone_url>",
			Username:         "<username>",
			Password:         "<password>",
			Domain:           "<domain>",
			TenantName:       "<tenant_name>",
		},
		CryptoConfig: config.CryptoConfig{
			EncryptionKey: key,
		},
	}
	connection := NewTestConnection()
	openstack, err := openstack2.NewClient(cfg)

	ctx := context.Background()
	if err != nil {
		t.Fatal("Failed to init openstack connection: ", err)
	}
	provisioner := NewMinecraftProvisioner(connection.Db, openstack, cfg.CryptoConfig.EncryptionKey)

	server, err := provisioner.NewGameServer(ctx, "Dev Server", types.Medium, types.Ubuntu24_04)
	if err != nil {
		t.Fatal("Failed to create server: ", err)
	}

	t.Log("Created server IP: ", server.Address)

	err = provisioner.DeleteGameServer(ctx, *server)
	t.Log("Deleted server IP: ", server.Address)
}

func MarshalED25519PrivateKey(key ed25519.PrivateKey) []byte {
	// Add our key header (followed by a null byte)
	magic := append([]byte("openssh-key-v1"), 0)

	var w struct {
		CipherName   string
		KdfName      string
		KdfOpts      string
		NumKeys      uint32
		PubKey       []byte
		PrivKeyBlock []byte
	}

	// Fill out the private key fields
	pk1 := struct {
		Check1  uint32
		Check2  uint32
		Keytype string
		Pub     []byte
		Priv    []byte
		Comment string
		Pad     []byte `ssh:"rest"`
	}{}

	// Set our check ints
	ci := rand.Uint32()
	pk1.Check1 = ci
	pk1.Check2 = ci

	// Set our key type
	pk1.Keytype = ssh.KeyAlgoED25519

	// Add the pubkey to the optionally-encrypted block
	pk, ok := key.Public().(ed25519.PublicKey)
	if !ok {
		//fmt.Fprintln(os.Stderr, "ed25519.PublicKey type assertion failed on an ed25519 public key. This should never ever happen.")
		return nil
	}
	pubKey := []byte(pk)
	pk1.Pub = pubKey

	// Add our private key
	pk1.Priv = []byte(key)

	// Might be useful to put something in here at some point
	pk1.Comment = ""

	// Add some padding to match the encryption block size within PrivKeyBlock (without Pad field)
	// 8 doesn't match the documentation, but that's what ssh-keygen uses for unencrypted keys. *shrug*
	bs := 8
	blockLen := len(ssh.Marshal(pk1))
	padLen := (bs - (blockLen % bs)) % bs
	pk1.Pad = make([]byte, padLen)

	// Padding is a sequence of bytes like: 1, 2, 3...
	for i := 0; i < padLen; i++ {
		pk1.Pad[i] = byte(i + 1)
	}

	// Generate the pubkey prefix "\0\0\0\nssh-ed25519\0\0\0 "
	prefix := []byte{0x0, 0x0, 0x0, 0x0b}
	prefix = append(prefix, []byte(ssh.KeyAlgoED25519)...)
	prefix = append(prefix, []byte{0x0, 0x0, 0x0, 0x20}...)

	// Only going to support unencrypted keys for now
	w.CipherName = "none"
	w.KdfName = "none"
	w.KdfOpts = ""
	w.NumKeys = 1
	w.PubKey = append(prefix, pubKey...)
	w.PrivKeyBlock = ssh.Marshal(pk1)

	magic = append(magic, ssh.Marshal(w)...)

	return magic
}
