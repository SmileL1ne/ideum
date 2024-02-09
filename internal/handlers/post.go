package handlers

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
)

/*
	TODO:
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

	comments, err := r.service.Comment.GetAllCommentsForPost(id)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidPostId):
			r.notFound(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	// Add comments to received post
	post.Comments = *comments

	// Create template struct with page related data
	data := r.newTemplateData(req)
	data.Models.Post = post
	data.Models.Comments = *comments
	data.Form = entity.CommentCreateForm{}

	r.render(w, req, http.StatusOK, "view.html", data)
}

func (r *routes) postCreate(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodPost:
		r.postCreatePost(w, req)
		return
	case req.Method != http.MethodGet:
		r.methodNotAllowed(w)
		return
	}

	data := r.newTemplateData(req)
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
	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	p := entity.PostCreateForm{Title: title, Content: content}

	id, err := r.service.Post.SavePost(&p, userID)
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

	// Add flash message to context
	r.sesm.Put(req.Context(), "flash", "Post successfully created!")

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}
