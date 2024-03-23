package entity

import "time"

type Request struct {
	ID        int
	UserID    int
	Username  string // not in db
	CreatedAt time.Time
}
