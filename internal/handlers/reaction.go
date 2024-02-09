package handlers

import (
	"fmt"
	"net/http"
)

func (r *routes) like(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodDelete:
		r.unLike(w, req)
		return
	case req.Method != http.MethodPatch:
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

	err := r.service.Reaction.Save(true, postID, userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%s", postID), http.StatusOK)
}

func (r *routes) unLike(w http.ResponseWriter, req *http.Request) {

}

func (r *routes) dislike(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodDelete:
		r.unDislike(w, req)
		return
	case req.Method != http.MethodPatch:
		r.methodNotAllowed(w)
		return
	}

}

func (r *routes) unDislike(w http.ResponseWriter, req *http.Request) {

}
