package http

import (
	"net/http"
)

/*
	TODO:
	- Implement all handlers
	- For post create add method check - if post -> redirect to postCreatePost handler
*/

func (r *routes) postView(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Display particular post"))
}

func (r *routes) postCreate(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Display new post creation page"))
}

func (r *routes) postCreatePost(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Create a new post"))
}
