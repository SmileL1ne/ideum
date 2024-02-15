package main

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strings"
)

/*
	TODO:
	- Find way to merge these 2 routes into 1, because they have pretty much
	same body except for 'like' it is true and for 'dislike' - false
	- Before liking post or comment check if they exist
*/

func (r *routes) postLike(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
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

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	err = r.service.Reaction.AddOrDeletePost(true, postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			r.badRequest(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}

func (r *routes) postDislike(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
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

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	err = r.service.Reaction.AddOrDeletePost(false, postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			r.badRequest(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}

func (r *routes) commentLike(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	commentID, err := r.getIdFromPath(req, 6)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidURLPath) {
			r.notFound(w)
			return
		}
		r.badRequest(w)
		return
	}

	urls := strings.Split(req.URL.Path, "/")
	postID := urls[len(urls)-2]

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	err = r.service.Reaction.AddOrDeleteComment(true, commentID, userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%s", postID), http.StatusSeeOther)
}

func (r *routes) commentDislike(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	commentID, err := r.getIdFromPath(req, 6)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidURLPath) {
			r.notFound(w)
			return
		}
		r.badRequest(w)
		return
	}

	urls := strings.Split(req.URL.Path, "/")
	postID := urls[len(urls)-2]

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	err = r.service.Reaction.AddOrDeleteComment(false, commentID, userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%s", postID), http.StatusSeeOther)
}
