package types

import "time"

// Backup represents a world backup with an associated server
type Backup struct {
	Id        int64
	ServerId  int64 `db:"server_id"`
	Timestamp time.Time
	Size      int
}
