package entity

// TODO: Add database related metatags

type Post struct {
	Title    string
	Content  string
	Likes    int
	Dislikes int
	UserID   int // Foreign key
}
