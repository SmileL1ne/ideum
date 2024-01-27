package http

import (
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strings"
)

/*
	TODO:
	- Check comments in methods
	- For postCreateForm - add userID for foreign key (when creating links between tables)
	- For post create add method check - if post -> redirect to postCreatePost handler
*/

func (r *routes) postView(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	urlParts := strings.Split(req.URL.Path, "/")
	if len(urlParts) != 4 {
		r.notFound(w)
		return
	}

	id := urlParts[len(urlParts)-1]

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
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	data := r.newTemplateData(req) // Empty template data (the original struct has no fields)
	data.Form = entity.PostCreateForm{}

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

	p := entity.PostCreateForm{Title: title, Content: content}

	id, status, err := r.service.Post.SavePost(&p)
	if status != http.StatusOK {
		if status == http.StatusUnprocessableEntity {
			data := r.newTemplateData(req)
			data.Form = p
			r.render(w, req, http.StatusUnprocessableEntity, "create.html", data)
		} else {
			r.serverError(w, req, err)
		}

		return
	}

	r.sesm.Put(req.Context(), "flash", "Post successfully created!")

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}
