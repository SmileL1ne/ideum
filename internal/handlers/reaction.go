package handlers

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strings"
)

func (r *Routes) postReaction(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	reaction := req.URL.Query().Get("reaction")
	if reaction == "" {
		r.logger.Print("postReaction: invalid query parameter - reaction")
		r.badRequest(w)
		return
	}

	postID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("postReaction: invalid url path")
		r.notFound(w)
		return
	}

	isPostExists, err := r.service.Post.ExistsPost(postID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if !isPostExists {
		r.logger.Printf("postReaction: no post with id - %d", postID)
		r.notFound(w)
		return
	}

	err = r.service.Reaction.SetPostReaction(reaction, postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			r.logger.Printf("postReaction: invalid query parameter - '%s'", reaction)
			r.badRequest(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}

func (r *Routes) commentReaction(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	reaction := req.URL.Query().Get("reaction")
	if reaction == "" {
		r.logger.Print("commentReaction: invalid query parameter - reaction")
		r.badRequest(w)
		return
	}

	commentID, ok := getIdFromPath(req, 6)
	if !ok {
		r.logger.Print("commentReaction: invalid url path")
		r.notFound(w)
		return
	}

	isCommentExist, err := r.service.Comment.ExistsComment(commentID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if !isCommentExist {
		r.logger.Printf("commentReaction: no comment with id - %d", commentID)
		r.notFound(w)
		return
	}

	urls := strings.Split(req.URL.Path, "/")
	postID, isValid := getValidID(urls[len(urls)-2])
	if !isValid {
		r.logger.Print("commentReaction: invalid postID")
		r.badRequest(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	err = r.service.Reaction.SetCommentReaction(reaction, commentID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			r.logger.Printf("commentReaction: invalid query parameter - '%s'", reaction)
			r.badRequest(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}
