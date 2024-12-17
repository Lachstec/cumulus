package db

import (
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"net"
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

	server := types.Server{
		Id:               0,
		Name:             "Test Server",
		Address:          net.ParseIP("192.168.1.1"),
		Status:           types.Stopped,
		Port:             1337,
		Memory:           2056,
		Game:             "Minecraft",
		GameVersion:      "1.0.0",
		GameMode:         types.Survival,
		Difficulty:       types.Normal,
		WhitelistEnabled: false,
		PlayersMax:       2,
		SSHKey:           []byte("sample ssh key"),
	}

	id, err := serverStore.Add(server)
	if err != nil {
		t.Fatalf("unable to save server: %s", err)
	}

	inserted, err := serverStore.GetById(id)
	if err != nil {
		t.Fatalf("unable to get inserted server from database: %s", err)
	}

	cmp.Equal(server, inserted)

	backup := types.Backup{
		Id:        0,
		ServerId:  inserted.Id,
		Timestamp: time.Now(),
		Size:      4096,
	}

	backupId, err := backupStore.Add(backup)
	if err != nil {
		t.Fatalf("unable to save backup: %s", err)
	}

	insertedBackup, err := backupStore.GetById(backupId)

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

	t.Log(inserted.Id)
	t.Log(updated.Id)

	if updated.Name != "This is a new name!" {
		t.Fatalf("expected name to be 'This is a new name!', got %s", updated.Name)
	}

	find, err := serverStore.Find(func(s types.Server) bool { return s.Id == updated.Id })
	if err != nil {
		t.Fatalf("unable to find server: %s", err)
	}

	cmp.Equal(find[0], updated)

	err = serverStore.Delete(find[0])
	if err != nil {
		t.Fatalf("unable to delete server: %s", err)
	}
}
