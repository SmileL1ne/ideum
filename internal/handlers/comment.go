package handlers

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strings"
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

	postID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("commentCreate: invalid url path")
		r.notFound(w)
		return
	}

	isPostExists, err := r.services.Post.ExistsPost(postID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if !isPostExists {
		r.logger.Printf("commentCreate: no post with id - %d", postID)
		r.notFound(w)
		return
	}

	content := req.PostForm.Get("commentContent")
	comment := &entity.CommentCreateForm{
		Content: content,
	}

	err = r.services.Comment.SaveComment(comment, postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			r.logger.Print("commentCreate: invalid form fill")
			w.WriteHeader(http.StatusBadRequest)
			msg := getErrorMessage(&comment.Validator)
			fmt.Fprint(w, strings.TrimSpace(msg))
		default:
			r.serverError(w, req, err)
		}
		return
	}

	redirectURL := fmt.Sprintf("/post/view/%d", postID)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, redirectURL)
}

func (r *Routes) commentDelete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	userRole := r.sesm.GetUserRole(req.Context())
	if userRole == entity.MODERATOR || userRole == entity.ADMIN {
		r.commentDeletePrivileged(w, req)
		return
	} else if userRole != entity.USER {
		r.unauthorized(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	commentID, ok := getIdFromPath(req, 5)
	if !ok {
		r.logger.Print("commentDelete: invalid url path")
		r.notFound(w)
		return
	}

	err := r.services.Comment.DeleteComment(commentID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrCommentNotFound):
			r.notFound(w)
		case errors.Is(err, entity.ErrForbiddenAccess):
			r.forbidden(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Routes) commentDeletePrivileged(w http.ResponseWriter, req *http.Request) {
	userRole := r.sesm.GetUserRole(req.Context())

	commentID, ok := getIdFromPath(req, 5)
	if !ok {
		r.logger.Print("commentDeletePrivileged: invalid url path")
		r.notFound(w)
		return
	}

	err := r.services.Comment.DeleteCommentPrivileged(commentID, userRole)
	if err != nil {
		if errors.Is(err, entity.ErrCommentNotFound) {
			r.notFound(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Routes) commentReport(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}
	if r.sesm.GetUserRole(req.Context()) != entity.MODERATOR {
		r.forbidden(w)
		return
	}

	commentID, ok := getIdFromPath(req, 5)
	if !ok {
		r.logger.Print("commentReport: invalid url path")
		r.notFound(w)
		return
	}
	exists, err := r.services.Comment.ExistsComment(commentID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if !exists {
		r.notFound(w)
		return
	}

	urls := strings.Split(req.URL.Path, "/")
	postID, isValid := getValidID(urls[len(urls)-2])
	if !isValid {
		r.logger.Print("commentReport: invalid postID")
		r.badRequest(w)
		return
	}

	message := req.PostFormValue("message")
	userID := r.sesm.GetUserID(req.Context())

	notification := entity.Notification{
		Type:       entity.REPORT,
		Content:    message,
		SourceID:   postID,
		SourceType: entity.COMMENT,
		UserFrom:   userID,
	}

	err = r.services.User.SendNotification(notification)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, req.URL.Path, http.StatusOK)
}
