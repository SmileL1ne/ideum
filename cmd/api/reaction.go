package main

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"log"
	"net/http"
	"strings"
)

func (r *routes) postReaction(w http.ResponseWriter, req *http.Request) {
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
		log.Print("postReaction: invalid query parameter - reaction")
		r.badRequest(w)
		return
	}

	postID, err := getIdFromPath(req, 4)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			log.Print("postReaction: invalid url path")
			r.notFound(w)
		case errors.Is(err, entity.ErrInvalidPathID):
			log.Print("postReaction: invalid id in request path")
			r.badRequest(w)
		}
		return
	}

	isPostExists, err := r.service.Post.ExistsPost(postID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if !isPostExists {
		log.Printf("postReaction: no post with id - %d", postID)
		r.notFound(w)
		return
	}

	err = r.service.Reaction.SetPostReaction(reaction, postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			log.Printf("postReaction: invalid query parameter - '%s'", reaction)
			r.badRequest(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}

func (r *routes) commentReaction(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	reaction := req.URL.Query().Get("reaction")
	if reaction == "" {
		log.Print("postReaction: invalid query parameter - reaction")
		r.badRequest(w)
		return
	}

	commentID, err := getIdFromPath(req, 6)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			log.Print("commentReaction: invalid url path")
			r.notFound(w)
		case errors.Is(err, entity.ErrInvalidPathID):
			log.Print("commentReaction: invalid id in request path")
			r.badRequest(w)
		}
		return
	}

	isCommentExists, err := r.service.Comment.ExistsComment(commentID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if !isCommentExists {
		log.Printf("commentReaction: no comment with id - %d", commentID)
		r.notFound(w)
		return
	}

	urls := strings.Split(req.URL.Path, "/")
	postID, isValid := getValidID(urls[len(urls)-2])
	if !isValid {
		log.Print("commentReaction: invalid postID")
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
			log.Printf("postReaction: invalid query parameter - '%s'", reaction)
			r.badRequest(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	http.Redirect(w, req, fmt.Sprintf("/post/view/%d", postID), http.StatusSeeOther)
}
