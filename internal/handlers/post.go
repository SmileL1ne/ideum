package handlers

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strings"
)

func (r *Routes) postView(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	username, tags, err := r.getBaseInfo(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	postID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("postView: invalid url path")
		r.notFound(w)
		return
	}

	post, err := r.service.Post.GetPost(postID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidPostID):
			r.logger.Print("postView: invalid post id")
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

	data := r.newTemplateData(req)
	data.Models.Post = post
	data.Models.Comments = *comments
	data.Models.Tags = *tags
	data.Username = username

	r.render(w, req, http.StatusOK, "view.html", data)
}

func (r *Routes) postCreate(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodPost:
		r.postCreatePost(w, req)
		return
	case req.Method != http.MethodGet:
		r.methodNotAllowed(w)
		return
	}

	username, tags, err := r.getBaseInfo(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Tags = *tags
	data.Username = username

	r.render(w, req, http.StatusOK, "create.html", data)
}

func (r *Routes) postCreatePost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseForm(); err != nil {
		r.logger.Print("postCreatePost: invalid form fill (parse error)")
		r.badRequest(w)
		return
	}

	form := req.PostForm

	title := strings.TrimSpace(form.Get("title"))
	content := strings.TrimSpace(form.Get("content"))

	// Get all selected tags id
	tags := form["tags"]
	if len(tags) == 0 {
		r.logger.Print("postCreatePost: no tags selected")
		w.WriteHeader(http.StatusBadRequest)
		msg := "tags: At least one tag should be selected"
		fmt.Fprint(w, strings.TrimSpace(msg))
		return
	}

	areTagsExist, err := r.service.Tag.AreTagsExist(tags)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			r.logger.Print("postCreatePost: invalid post tags")

			w.WriteHeader(http.StatusBadRequest)
			msg := "invalid post tags"
			fmt.Fprint(w, strings.TrimSpace(msg))
		default:
			r.serverError(w, req, err)
		}
		return
	}
	if !areTagsExist {
		r.logger.Print("postCreatePost: post tags don't exist")

		w.WriteHeader(http.StatusBadRequest)
		msg := "invalid post tags (don't exist)"
		fmt.Fprint(w, strings.TrimSpace(msg))
		return
	}

	// Get userID from request
	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	p := entity.PostCreateForm{
		Title:   title,
		Content: content,
	}

	id, err := r.service.Post.SavePost(&p, userID, tags)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			r.logger.Print("postCreatePost: invalid form fill")

			w.WriteHeader(http.StatusBadRequest)
			msg := getErrorMessage(&p.Validator)
			fmt.Fprint(w, strings.TrimSpace(msg))
		default:
			r.serverError(w, req, err)
		}
		return
	}

	redirectURL := fmt.Sprintf("/post/view/%d", id)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, redirectURL)
}

func (r *Routes) postsPersonal(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	username, tags, err := r.getBaseInfo(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if username == "" { // This should never happen
		r.unauthorized(w)
		return
	}

	userPosts, err := r.service.Post.GetAllPostsByUserId(userID)
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

func (r *Routes) postsReacted(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	username, tags, err := r.getBaseInfo(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if username == "" {
		r.unauthorized(w)
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
