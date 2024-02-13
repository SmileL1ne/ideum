package handlers

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
)

func (r *routes) commentCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
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

	content := req.PostForm.Get("commentContent")
	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	comment := &entity.CommentCreateForm{
		Content: content,
	}

	err = r.service.Comment.SaveComment(comment, postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			http.Redirect(w, req, fmt.Sprintf("/post/view/%d", postID), http.StatusBadRequest)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	// Add flash message to request's context
	r.sesm.Put(req.Context(), "flash", "Successfully added your comment!")

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}
