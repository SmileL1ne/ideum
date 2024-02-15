package main

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"log"
	"net/http"
)

func (r *routes) postView(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	postID, err := r.getIdFromPath(req, 4)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidURLPath) {
			r.notFound(w)
			return
		}
		r.badRequest(w)
		return
	}

	post, err := r.service.Post.GetPost(postID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	comments, err := r.service.Comment.GetAllCommentsForPost(postID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	tags, err := r.service.Tag.GetAllTags()
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	postTags, err := r.service.Tag.GetAllTagsForPost(postID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Post = post
	data.Models.Post.PostTags = *postTags
	data.Models.Comments = *comments
	data.Models.Tags = *tags

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

	tags, err := r.service.Tag.GetAllTags()
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Tags = *tags

	r.render(w, req, http.StatusOK, "create.html", data)
}

func (r *routes) postCreatePost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseForm(); err != nil {
		log.Print("postCreatePost: invalid form fill (parse error)")
		r.badRequest(w)
		return
	}

	form := req.PostForm

	title := form.Get("title")
	content := form.Get("content")

	// Get all selected tags id
	tags := form["tags"]

	// Get userID from request
	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	p := entity.PostCreateForm{Title: title, Content: content}

	id, err := r.service.Post.SavePost(&p, userID, tags)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			log.Print("postCreatePost: invalid form fill")
			http.Redirect(w, req, "/post/create", http.StatusBadRequest)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	redirectURL := fmt.Sprintf("/post/view/%d", id)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, redirectURL)

}
