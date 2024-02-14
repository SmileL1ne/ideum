package main

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"log"
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

	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	postID, err := getIdFromPath(req, 4)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			log.Print("commentCreate: invalid url path")
			r.notFound(w)
		case errors.Is(err, entity.ErrInvalidPathID):
			log.Print("commentCreate: invalid id in request path")
			r.badRequest(w)
		}
		return
	}

	content := req.PostForm.Get("commentContent")
	comment := &entity.CommentCreateForm{
		Content: content,
	}

	err = r.service.Comment.SaveComment(comment, postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			log.Print("commentCreate: invalid form fill")
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
