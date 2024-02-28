package entity

import "time"

// TagEntity is returned from repo and from service
type TagEntity struct {
	ID        int
	Name      string
	CreatedAt time.Time
}
