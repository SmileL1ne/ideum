package entity

import "time"

// TODO: Add database related metatags

type Post struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	// UserID  int // Foreign key
}
