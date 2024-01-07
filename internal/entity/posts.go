package entity

// TODO: Add database related metatags

type Post struct {
	Title    string
	Content  string
	Likes    int
	Dislikes int
	userID   int // Foreign key
}
