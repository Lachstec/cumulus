package services

import (
	"github.com/Lachstec/mc-hosting/internal/db"
	"github.com/Lachstec/mc-hosting/internal/types"
)

type FloatingIPService struct {
	store db.Store[types.FloatingIP]
}

func NewFloatingIPService(store db.Store[types.FloatingIP]) *FloatingIPService {
	return &FloatingIPService{
		store: store,
	}
}

func (c *FloatingIPService) ReadIpByServerID(serverid int64) ([]*types.FloatingIP, error) {
	floatingip, err := c.store.Find(func(s *types.FloatingIP) bool { return s.Id == serverid })
	if err != nil {
		return nil, err
	}
	return floatingip, nil
}
