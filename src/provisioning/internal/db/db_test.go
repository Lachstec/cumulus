package db

import (
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type TestConnection struct {
	Db *sqlx.DB
}

func NewTestConnection() *TestConnection {
	conn, err := sqlx.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		panic("unable to connect to testing database")
	}
	migrator := NewMigrator(conn)
	err = migrator.Migrate("../../migrations")
	if err != nil {
		panic(err)
	}

	return &TestConnection{conn}
}

func TestServerStore(t *testing.T) {
	db := NewTestConnection()
	serverStore := NewServerStore(db.Db)
	backupStore := NewServerBackupStore(db.Db)
	userStore := NewUserStore(db.Db)

	user := &types.User{
		Sub:   "Bla",
		Name:  "Testuser",
		Class: string(types.Admin),
	}
	keyStore := NewKeyStore(db.Db)
	ipStore := NewIPStore(db.Db)

	key := types.Key{
		Name:       "key1",
		PrivateKey: []byte{1, 2, 3},
		PublicKey:  []byte{4, 5, 6},
	}

	kid, err := keyStore.Add(&key)
	if err != nil {
		t.Fatal(err)
	}

	uid, err := userStore.Add(user)
	if err != nil {
		t.Fatal(err)
	}

	addr := types.FloatingIP{
		OpenstackId: "696d03a0-c5c3-4eba-878f-b0aca2d84cd0",
		Ip:          "192.168.0.1",
	}

	iid, err := ipStore.Add(addr)
	if err != nil {
		t.Fatal(err)
	}

	server := &types.Server{
		UserID:           uid,
		OpenstackId:      "31e0683c-5455-4510-b3ba-3c02241a3eff",
		Name:             "Test Server",
		Address:          iid,
		Status:           types.Stopped,
		Port:             1337,
		Flavour:          types.Flavours[2].ID,
		Image:            "Ubuntu 20.04",
		Game:             "Minecraft",
		GameVersion:      "1.0.0",
		GameMode:         types.Survival,
		Difficulty:       types.Normal,
		WhitelistEnabled: false,
		PvPEnabled:       true,
		PlayersMax:       2,
		SSHKey:           kid,
	}

	ID, err := serverStore.Add(server)
	if err != nil {
		t.Fatalf("unable to save server: %s", err)
	}

	inserted, err := serverStore.GetById(ID)
	if err != nil {
		t.Fatalf("unable to get inserted server from database: %s", err)
	}

	cmp.Equal(server, inserted)

	backup := &types.Backup{
		ID:          0,
		OpenstackID: "31e0683c-5455-4510-b3ba-3c02241a3eff",
		ServerID:    inserted.ID,
		Timestamp:   time.Now(),
		Size:        4096,
	}

	backupID, err := backupStore.Add(backup)
	if err != nil {
		t.Fatalf("unable to save backup: %s", err)
	}

	insertedBackup, err := backupStore.GetById(backupID)

	if err != nil {
		t.Fatalf("unable to get inserted backup from database: %s", err)
	}

	cmp.Equal(insertedBackup, backup)
	updatedBackup, err := backupStore.Update(insertedBackup)

	if err != nil {
		t.Fatalf("unable to update backup: %s", err)
	}

	err = backupStore.Delete(updatedBackup)
	if err != nil {
		t.Fatalf("unable to delete backup: %s", err)
	}

	inserted.Name = "This is a new name!"

	updated, err := serverStore.Update(inserted)

	if err != nil {
		t.Fatalf("unable to update server: %s", err)
	}

	t.Log(inserted.ID)
	t.Log(updated.ID)

	if updated.Name != "This is a new name!" {
		t.Fatalf("expected name to be 'This is a new name!', got %s", updated.Name)
	}

	find, err := serverStore.Find(func(s *types.Server) bool { return s.ID == updated.ID })
	if err != nil {
		t.Fatalf("unable to find server: %s", err)
	}

	cmp.Equal(find[0], updated)

	err = serverStore.Delete(find[0])
	if err != nil {
		t.Fatalf("unable to delete server: %s", err)
	}
}

func TestKeyStore(t *testing.T) {
	db := NewTestConnection()
	keyStore := NewKeyStore(db.Db)

	id, err := keyStore.Add(&types.Key{Name: "Test Key", PrivateKey: []byte("test key"), PublicKey: []byte("test pub key")})
	if err != nil {
		t.Fatalf("unable to save key: %s", err)
	}

	key, err := keyStore.GetById(id)
	if err != nil {
		t.Fatalf("unable to get key: %s", err)
	}

	assert.Equal(t, key.Name, "Test Key")
	assert.Equal(t, key.PrivateKey, []byte("test key"))
	assert.Equal(t, key.PublicKey, []byte("test pub key"))

	key.Name = "Another Key"
	updated, err := keyStore.Update(key)
	if err != nil {
		t.Fatalf("unable to update key: %s", err)
	}

	assert.Equal(t, updated.Name, "Another Key")

	err = keyStore.Delete(key)
	if err != nil {
		t.Fatalf("unable to delete key: %s", err)
	}

	keys, err := keyStore.Find(func(k *types.Key) bool { return k.Id == updated.Id })
	if err != nil {
		t.Fatalf("unable to find keys: %s", err)
	}

	if len(keys) != 0 {
		t.Fatalf("expected 0 keys, got %d", len(keys))
	}
}

func TestIpStore(t *testing.T) {
	db := NewTestConnection()
	ipStore := NewIPStore(db.Db)

	id, err := ipStore.Add(types.FloatingIP{OpenstackId: "696d03a0-c5c3-4eba-878f-b0aca2d84cd0", Ip: "192.168.0.1"})
	if err != nil {
		t.Fatalf("unable to save ip: %s", err)
	}

	ip, err := ipStore.GetById(id)
	if err != nil {
		t.Fatalf("unable to get ip: %s", err)
	}

	assert.Equal(t, ip.OpenstackId, "696d03a0-c5c3-4eba-878f-b0aca2d84cd0")
	assert.Equal(t, ip.Ip, "192.168.0.1")

	ip.OpenstackId = "307c6705-02b2-4144-94ad-2be61134f283"
	updated, err := ipStore.Update(ip)
	if err != nil {
		t.Fatalf("unable to update ip: %s", err)
	}

	assert.Equal(t, updated.OpenstackId, "307c6705-02b2-4144-94ad-2be61134f283")

	err = ipStore.Delete(ip)
	if err != nil {
		t.Fatalf("unable to delete ip: %s", err)
	}

	ips, err := ipStore.Find(func(i types.FloatingIP) bool { return i.Id == updated.Id })
	if err != nil {
		t.Fatalf("unable to find ips: %s", err)
	}

	if len(ips) != 0 {
		t.Fatalf("expected 0 ips, got %d", len(ips))
	}
}
