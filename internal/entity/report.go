package entity

import "time"

type Report struct {
	ID         int
	Reason     string
	UserFrom   int
	Username   string // not in db
	SourceID   int
	SourceType string
	CreatedAt  time.Time
}
