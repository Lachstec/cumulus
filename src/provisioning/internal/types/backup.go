package types

import "time"

// Backup represents a world backup with an associated server
type Backup struct {
	ID          int64
	OpenstackID string `db:"openstack_id"`
	ServerID    int64  `db:"server_id"`
	Timestamp   time.Time
	Size        int
}
