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

	username, err := r.getUsername(req)
	if err != nil && !errors.Is(err, entity.ErrInvalidUserID) {
		r.serverError(w, req, err)
		return
	}

	postID, err := getIdFromPath(req, 4)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			log.Print("postView: invalid url path")
			r.notFound(w)
		case errors.Is(err, entity.ErrInvalidPathID):
			log.Print("postView: invalid id in request path")
			r.badRequest(w)
		}
		return
	}

	post, err := r.service.Post.GetPost(postID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidPostID):
			log.Print("postView: invalid post id")
			r.notFound(w)
		default:
			r.serverError(w, req, err)
		}
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

	data := r.newTemplateData(req)
	data.Models.Post = post
	data.Models.Comments = *comments
	data.Models.Tags = *tags
	data.Username = username

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

	username, err := r.getUsername(req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidUserID):
			log.Print("postCreate: invalid user id")
			r.unauthorized(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	tags, err := r.service.Tag.GetAllTags()
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Tags = *tags
	data.Username = username

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

	// Add flash message to context
	r.sesm.Put(req.Context(), "flash", "Post successfully created!")

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", id), http.StatusSeeOther)
}

func (r *routes) postsPersonal(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	username, err := r.getUsername(req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidUserID):
			log.Print("postCreate: invalid user id")
			r.unauthorized(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	tags, err := r.service.Tag.GetAllTags()
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	userPosts, err := r.service.Post.GetAllPostsByUserID(userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Tags = *tags
	data.Models.Posts = *userPosts
	data.Username = username

	r.render(w, req, http.StatusOK, "home.html", data)
}

func (r *routes) postsReacted(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	username, err := r.getUsername(req)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidUserID):
			log.Print("postCreate: invalid user id")
			r.unauthorized(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	tags, err := r.service.Tag.GetAllTags()
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	reactedPosts, err := r.service.Post.GetAllPostsByUserReaction(userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Tags = *tags
	data.Models.Posts = *reactedPosts
	data.Username = username

	r.render(w, req, http.StatusOK, "home.html", data)
}