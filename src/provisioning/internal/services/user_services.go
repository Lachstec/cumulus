package services

import (
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/types"
)

type UserService struct {
	store db.Store[types.User]
}

func NewUserService(store db.Store[types.User]) *UserService {
	return &UserService{
		store: store,
	}
}

func (c *UserService) ReadAllUsers() ([]*types.User, error) {

	users, err := c.store.Find(func(s *types.User) bool { return true })
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *UserService) ReadUserByUserID(userid int64) ([]*types.User, error) {
	user, err := c.store.Find(func(s *types.User) bool { return s.ID == userid })
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (c *UserService) CreateUser(user *types.User) (int64, error) {
	userid, err := c.store.Add(user)
	if err != nil {
		return 0, err
	}
	return userid, nil
}

func (c *UserService) DeleteUser(user *types.User) (error) {
	err := c.store.Delete(user)
	if err != nil {
		return err
	}
	return nil
}

func (c *UserService) UpdateUser(user *types.User) (*types.User, error) {
	user, err := c.store.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
