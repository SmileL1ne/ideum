package handlers

import (
	"fmt"
	"net/http"
)

/*
	TODO:
	- Find way to merge these 2 routes into 1, because they have pretty much
	same body except for like it is true and for dislike - false
*/

func (r *routes) like(w http.ResponseWriter, req *http.Request) {
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

	err := r.service.Reaction.AddOrDelete(true, postID, userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%s", postID), http.StatusSeeOther)
}

func (r *routes) dislike(w http.ResponseWriter, req *http.Request) {
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

	err := r.service.Reaction.AddOrDelete(false, postID, userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%s", postID), http.StatusSeeOther)
}
