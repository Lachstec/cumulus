package services

import (
	"errors"
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockUserStore struct {
	users []types.User
}

func NewMockUserStore() db.Store[types.User] {
	return &MockUserStore{
		users: []types.User{
			{ID: 1, Sub: "sub1", Name: "Alice", Class: "Admin"},
			{ID: 2, Sub: "sub2", Name: "Bob", Class: "User"},
			{ID: 3, Sub: "sub3", Name: "Charlie", Class: "User"},
		},
	}
}

func (s *MockUserStore) GetById(id int64) (*types.User, error) {
	for i := range s.users {
		if s.users[i].ID == id {
			return &s.users[i], nil
		}
	}
	return nil, errors.New("not found")
}

func (s *MockUserStore) Add(user *types.User) (int64, error) {
	user.ID = int64(len(s.users) + 1)
	s.users = append(s.users, *user)
	return user.ID, nil
}

func (s *MockUserStore) Find(predicate db.Predicate[*types.User]) ([]*types.User, error) {
	var users []*types.User
	for i := range s.users {
		if predicate(&s.users[i]) {
			users = append(users, &s.users[i])
		}
	}
	return users, nil
}

func (s *MockUserStore) Delete(user *types.User) error {
	for i, u := range s.users {
		if u.ID == user.ID {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (s *MockUserStore) Update(user *types.User) (*types.User, error) {
	for i := range s.users {
		if s.users[i].ID == user.ID {
			s.users[i] = *user
			return &s.users[i], nil
		}
	}
	return nil, errors.New("user not found")
}

func TestReadAllUsers(t *testing.T) {
	service := NewUserService(NewMockUserStore())
	users, err := service.ReadAllUsers()
	if err != nil {
		t.Fatal("Unexpected error reading all users:", err)
	}

	assert.Equal(t, len(users), 3, "There should be 3 users")
}

func TestReadUser(t *testing.T) {
	service := NewUserService(NewMockUserStore())
	user, err := service.ReadUserByUserID(1)
	if err != nil {
		t.Fatal("Unexpected error reading user:", err)
	}

	if len(user) != 1 {
		t.Fatal("Expected one user")
	}

	assert.Equal(t, user[0].ID, int64(1), "The id should be 1")
}

func TestDeleteUser(t *testing.T) {
	service := NewUserService(NewMockUserStore())
	user, err := service.ReadUserByUserID(1)
	if err != nil {
		t.Fatal("Unexpected error reading user:", err)
	}

	if len(user) != 1 {
		t.Fatal("Expected one user")
	}

	err = service.DeleteUser(user[0])
	if err != nil {
		t.Fatal("Unexpected error deleting user:", err)
	}

	users, err := service.ReadAllUsers()
	if err != nil {
		t.Fatal("Unexpected error reading all users:", err)
	}

	assert.Equal(t, 2, len(users), "There should be two users")
}

func TestUpdateUser(t *testing.T) {
	service := NewUserService(NewMockUserStore())
	user, err := service.ReadUserByUserID(1)
	if err != nil {
		t.Fatal("Unexpected error reading user:", err)
	}

	if len(user) != 1 {
		t.Fatal("Expected one user")
	}

	user[0].Name = "Soyjak"

	u, err := service.UpdateUser(user[0])
	if err != nil {
		t.Fatal("Unexpected error updating user:", err)
	}

	assert.Equal(t, u.Name, "Soyjak", "The name should be Soyjak")
}

func TestNewUser(t *testing.T) {
	service := NewUserService(NewMockUserStore())
	user := types.User{
		Sub:   "Subby",
		Name:  "Soyboy",
		Class: types.Admin.Value(),
	}

	_, err := service.CreateUser(&user)
	if err != nil {
		t.Fatal("Unexpected error creating user:", err)
	}

	users, err := service.ReadAllUsers()
	if err != nil {
		t.Fatal("Unexpected error reading all users:", err)
	}

	assert.Equal(t, 4, len(users), "There should be four users")
}
