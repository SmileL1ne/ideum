package entity

import "time"

type TagEntity struct {
	ID        int
	Name      string
	CreatedAt time.Time
}
