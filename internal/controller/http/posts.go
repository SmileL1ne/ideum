package http

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
)

/*
	TODO:
	- Check comments in methods
	- For postCreateForm - add userID for foreign key (when creating links between tables)
*/

func (r *routes) postView(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	id := r.getIdFromPath(req, 4)
	if id == "" {
		r.notFound(w)
		return
	}

	post, err := r.service.Post.GetPost(id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrNoRecord):
			r.notFound(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	commentsForPost, err := r.service.Comment.GetAllCommentsForPost(id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidPostId):
			r.notFound(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	data := r.newTemplateData(req)
	data.Post = post
	data.Comments = commentsForPost
	data.Form = entity.CommentCreateForm{}

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

	id, err := r.service.Post.SavePost(&p)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			data := r.newTemplateData(req)
			data.Form = p
			r.render(w, req, http.StatusUnprocessableEntity, "create.html", data)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	r.sesm.Put(req.Context(), "flash", "Post successfully created!")

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}
