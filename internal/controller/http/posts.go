package http

import (
	"fmt"
	"forum/internal/entity"
	"net/http"
)

/*
	TODO:
	- Check comments in methods
	- For postCreateForm - add userID for foreign key (when creating links between tables)
	- For post create add method check - if post -> redirect to postCreatePost handler
*/

func (r *routes) postView(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.notFound(w)
		return
	}

	id := req.URL.Query().Get("id")
	post, status, err := r.service.Post.GetPost(id)
	if status != http.StatusOK {
		r.identifyStatus(w, req, status, err)
		return
	}

	data := r.newTemplateData(req)
	data.Post = post

	r.render(w, req, http.StatusOK, "view.html", data)
}

func (r *routes) postCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.postCreatePost(w, req)
		return
	} else if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	data := r.newTemplateData(req) // Empty template data (the original struct has no fields)

	r.render(w, req, http.StatusOK, "create.html", data)
}

func (r *routes) postCreatePost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}

	form := req.PostForm

	title := form.Get("title")
	content := form.Get("content")
	// Add user id here (foreign key)

	id, status, err := r.service.Post.SavePost(entity.Post{Title: title, Content: content})
	if status != http.StatusOK {
		r.identifyStatus(w, req, status, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view?id=%d", id), http.StatusSeeOther)
}
