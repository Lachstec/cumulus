package services

import (
	"errors"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockStore struct {
	servers []types.Server
}

func NewMockStore() db.Store[types.Server] {
	return &MockStore{
		servers: []types.Server{
			{
				ID:               1,
				UserID:           1337,
				OpenstackId:      "asdasdwasdwadasdsa",
				Name:             "Testserver1",
				Address:          2,
				Status:           types.Running,
				Port:             25560,
				Flavour:          types.Flavours[0].ID,
				Image:            types.Alpine3_20_3,
				Game:             "Minecraft",
				GameVersion:      "1.0.0",
				GameMode:         types.Survival,
				Difficulty:       types.Peaceful,
				WhitelistEnabled: true,
				PvPEnabled:       false,
				PlayersMax:       20,
				SSHKey:           3,
			},
			{
				ID:               2,
				UserID:           1337,
				OpenstackId:      "asdasdwasdwadasdsa",
				Name:             "Testserver2",
				Address:          2,
				Status:           types.Running,
				Port:             25560,
				Flavour:          types.Flavours[0].ID,
				Image:            types.Alpine3_20_3,
				Game:             "Minecraft",
				GameVersion:      "1.0.0",
				GameMode:         types.Survival,
				Difficulty:       types.Peaceful,
				WhitelistEnabled: true,
				PvPEnabled:       false,
				PlayersMax:       20,
				SSHKey:           3,
			},
			{
				ID:               3,
				UserID:           1337,
				OpenstackId:      "asdasdwasdwadasdsa",
				Name:             "Testserver3",
				Address:          2,
				Status:           types.Running,
				Port:             25560,
				Flavour:          types.Flavours[0].ID,
				Image:            types.Alpine3_20_3,
				Game:             "Minecraft",
				GameVersion:      "1.0.0",
				GameMode:         types.Survival,
				Difficulty:       types.Peaceful,
				WhitelistEnabled: true,
				PvPEnabled:       false,
				PlayersMax:       20,
				SSHKey:           3,
			},
		},
	}
}

func (s *MockStore) GetById(id int64) (*types.Server, error) {
	for i := range s.servers {
		if s.servers[i].ID == id {
			return &s.servers[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (s *MockStore) Add(server *types.Server) (int64, error) {
	s.servers = append(s.servers, *server)
	server.ID = int64(len(s.servers))
	return server.ID, nil
}

func (s *MockStore) Find(predicate db.Predicate[*types.Server]) ([]*types.Server, error) {
	var servers []*types.Server
	for i := range s.servers {
		if predicate(&s.servers[i]) {
			servers = append(servers, &s.servers[i])
		}
	}
	return servers, nil
}

func (s *MockStore) Delete(server *types.Server) error {
	for i, serv := range s.servers {
		if serv.ID == server.ID {
			s.servers = append(s.servers[:i], s.servers[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (s *MockStore) Update(server *types.Server) (*types.Server, error) {
	for i := range s.servers {
		if s.servers[i].ID == server.ID {
			s.servers[i] = *server // ✅ Update the existing record
			return &s.servers[i], nil
		}
	}
	return nil, errors.New("server not found")
}

func TestReadAllServers(t *testing.T) {
	t.Parallel()
	service := NewServerService(NewMockStore())
	servers, err := service.ReadAllServers()
	if err != nil {
		t.Fatalf("Unexpected error reading servers from database: %v", err)
	}

	assert.Equal(t, len(servers), 3, "There should be two servers in the database")
}

func TestReadServerByServerId(t *testing.T) {
	t.Parallel()
	service := NewServerService(NewMockStore())
	servers, err := service.ReadServerByServerID(1)
	if err != nil {
		t.Fatalf("Unexpected error reading servers from database: %v", err)
	}

	assert.Equal(t, len(servers), 1, "There should be one server in the Result")
	assert.Equal(t, servers[0].ID, int64(1), "Result should contain the server id 1")
}

func TestCreateServer(t *testing.T) {
	t.Parallel()
	service := NewServerService(NewMockStore())
	server := types.Server{
		ID:               3,
		UserID:           1337,
		OpenstackId:      "asdasdwasdwadasdsa",
		Name:             "Testserver3",
		Address:          2,
		Status:           types.Running,
		Port:             25560,
		Flavour:          types.Flavours[0].ID,
		Image:            types.Alpine3_20_3,
		Game:             "Minecraft",
		GameVersion:      "1.0.0",
		GameMode:         types.Survival,
		Difficulty:       types.Peaceful,
		WhitelistEnabled: true,
		PvPEnabled:       false,
		PlayersMax:       20,
		SSHKey:           3,
	}

	id, err := service.CreateServer(&server)

	if err != nil {
		t.Fatalf("Unexpected error creating server: %v", err)
	}
	assert.Equal(t, int64(4), int64(4), "There should be four servers in the database")
	assert.Equal(t, id, int64(4), "ID should be four")
}

func TestDeleteServer(t *testing.T) {
	t.Parallel()
	service := NewServerService(NewMockStore())

	server, err := service.ReadServerByServerID(1)
	if err != nil {
		t.Fatalf("Unexpected error reading servers from database: %v", err)
	}

	if len(server) == 0 {
		t.Fatal("There should be some servers in the database")
	}

	err = service.DeleteServer(server[0])
	if err != nil {
		t.Fatalf("Unexpected error deleting server: %v", err)
	}

	all, err := service.ReadAllServers()
	if err != nil {
		t.Fatalf("Unexpected error reading servers from database: %v", err)
	}

	assert.Equal(t, 2, len(all), "There should be two servers in the database")
}

func TestUpdateServer(t *testing.T) {
	t.Parallel()
	service := NewServerService(NewMockStore())
	server, err := service.ReadServerByServerID(1)
	if err != nil {
		t.Fatalf("Unexpected error reading servers from database: %v", err)
	}

	if len(server) == 0 {
		t.Fatal("There should be some servers in the database")
	}

	server[0].Name = "Soyserver"
	newServer, err := service.UpdateServer(server[0])
	if err != nil {
		t.Fatalf("Unexpected error updating server: %v", err)
	}

	assert.Equal(t, server[0].Name, newServer.Name, "Name did not get updated")
}
