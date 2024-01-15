package http

import (
	"fmt"
	"forum/internal/entity"
	"net/http"
)

/*
	TODO:
	- for postCreateForm - add userID for foreign key (when creating links between tables)
	- For post create add method check - if post -> redirect to postCreatePost handler
*/

func (r *routes) postView(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Display particular post"))
}

func (r *routes) postCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.postCreatePost(w, req)
		return
	} else if req.Method != http.MethodGet {
		r.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Display new post creation page"))
}

func (r *routes) postCreatePost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Black matter bubbles"
	content := "I suggest there should be black matter bubbles..."
	// userId := 10

	id, err := r.s.Post.SavePost(entity.Post{Title: title, Content: content})
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
}
