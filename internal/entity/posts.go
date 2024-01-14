package entity

// TODO: Add database related metatags

type Post struct {
	Title   string
	Content string
	UserID  int // Foreign key
}
