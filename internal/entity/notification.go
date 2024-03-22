package entity

import "time"

type Notification struct {
	ID         int
	Type       string
	Content    string
	SourceID   int    // source is id of the source of action (postID)
	SourceType string // source type can be either 'post' or 'comment'
	UserFrom   int
	UserTo     int
	Username   string // Not in db
	CreatedAt  time.Time
}
