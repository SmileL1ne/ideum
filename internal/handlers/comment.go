package handlers

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
)

func (r *Routes) commentCreate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	postID, err := getIdFromPath(req, 4)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			r.logger.Print("commentCreate: invalid url path")
			r.notFound(w)
		case errors.Is(err, entity.ErrInvalidPathID):
			r.logger.Print("commentCreate: invalid id in request path")
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
			r.logger.Print("commentCreate: invalid form fill")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, comment.FieldErrors["commentContent"])
		default:
			r.serverError(w, req, err)
		}
		return
	}

	redirectURL := fmt.Sprintf("/post/view/%d", postID)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, redirectURL)
}
