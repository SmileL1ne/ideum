package entity

type Notification struct {
	Type     string
	Content  string
	SourceID int // source is id of the source of action (postID or commentID)
	UserFrom int
	UserTo   int
}
