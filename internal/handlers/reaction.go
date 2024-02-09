package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

/*
	TODO:
	- Find way to merge these 2 routes into 1, because they have pretty much
	same body except for like it is true and for dislike - false
*/

func (r *routes) postLike(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	postID := r.getIdFromPath(req, 4)
	if postID == "" {
		r.notFound(w)
		return
	}

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	err := r.service.Reaction.AddOrDeletePost(true, postID, userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%s", postID), http.StatusSeeOther)
}

func (r *routes) postDislike(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	postID := r.getIdFromPath(req, 4)
	if postID == "" {
		r.notFound(w)
		return
	}

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	err := r.service.Reaction.AddOrDeletePost(false, postID, userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%s", postID), http.StatusSeeOther)
}

func (r *routes) commentLike(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	commentID := r.getIdFromPath(req, 6)
	if commentID == "" {
		r.notFound(w)
		return
	}

	urls := strings.Split(req.URL.Path, "/")
	postID := urls[len(urls)-2]

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	err := r.service.Reaction.AddOrDeleteComment(true, commentID, userID)
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

	commentID := r.getIdFromPath(req, 6)
	if commentID == "" {
		r.notFound(w)
		return
	}

	urls := strings.Split(req.URL.Path, "/")
	postID := urls[len(urls)-2]

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	err := r.service.Reaction.AddOrDeleteComment(false, commentID, userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%s", postID), http.StatusSeeOther)
}
