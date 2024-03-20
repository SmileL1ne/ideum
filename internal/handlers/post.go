package handlers

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strings"
)

func (r *Routes) postView(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	data, err := r.newTemplateData(req)
	if err != nil {
		if errors.Is(err, entity.ErrUnauthorized) {
			r.unauthorized(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	postID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("postView: invalid url path")
		r.notFound(w)
		return
	}

	post, err := r.services.Post.GetPost(postID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidPostID):
			r.logger.Print("postView: invalid post id")
			r.notFound(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	comments, err := r.services.Comment.GetAllCommentsForPost(postID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data.Models.Post = post
	data.Models.Post.Comments = *comments
	data.Models.Post.ImageName = post.ImageName

	r.render(w, req, http.StatusOK, "view.html", data)
}

func (r *Routes) postCreate(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodPost:
		r.postCreatePost(w, req)
		return
	case req.Method != http.MethodGet:
		r.methodNotAllowed(w)
		return
	}

	data, err := r.newTemplateData(req)
	if err != nil {
		if errors.Is(err, entity.ErrUnauthorized) {
			r.unauthorized(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	r.render(w, req, http.StatusOK, "create.html", data)
}

func (r *Routes) postCreatePost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseMultipartForm(10); err != nil {
		r.logger.Print("postCreatePost: invalid form fill (parse error)")
		r.badRequest(w)
		return
	}

	form := req.PostForm
	title := strings.TrimSpace(form.Get("title"))
	content := strings.TrimSpace(form.Get("content"))
	tags := form["tags"]

	// Take image file from file form
	var withImage bool = true
	file, fileHeader, imgErr := req.FormFile("image")
	if imgErr != nil {
		withImage = false
		r.logger.Print("postCreatePost: no file")
	}

	// Get userID from request's context
	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	p := entity.PostCreateForm{
		Title:      title,
		Content:    content,
		UserID:     userID,
		Tags:       tags,
		File:       file,
		FileHeader: fileHeader,
	}

	isPostValid, err := r.services.Post.CheckPostAttrs(&p, withImage)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidTags):
			r.logger.Print("postCreatePost: post tags don't exist")
			r.badRequest(w)
		default:
			r.serverError(w, req, err)
		}
		return
	}
	if !isPostValid {
		r.logger.Print("postCreatePost: invalid form fill")
		w.WriteHeader(http.StatusBadRequest)
		msg := getErrorMessage(&p.Validator)
		fmt.Fprint(w, strings.TrimSpace(msg))
		return
	}

	if withImage {
		imgName, err := r.services.Image.ProcessImage(file, fileHeader)
		if err != nil {
			r.serverError(w, req, err)
			return
		}
		p.ImageName = imgName
	}

	id, err := r.services.Post.SavePost(p)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	redirectURL := fmt.Sprintf("/post/view/%d", id)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, redirectURL)
}

func (r *Routes) postsPersonal(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	data, err := r.newTemplateData(req)
	if err != nil {
		if errors.Is(err, entity.ErrUnauthorized) {
			r.unauthorized(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	userPosts, err := r.services.Post.GetAllPostsByUserId(userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data.Models.Posts = *userPosts

	r.render(w, req, http.StatusOK, "home.html", data)
}

func (r *Routes) postsReacted(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	data, err := r.newTemplateData(req)
	if err != nil {
		if errors.Is(err, entity.ErrUnauthorized) {
			r.unauthorized(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	reactedPosts, err := r.services.Post.GetAllPostsByUserReaction(userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data.Models.Posts = *reactedPosts

	r.render(w, req, http.StatusOK, "home.html", data)
}

func (r *Routes) postsCommented(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID, data, err := r.getBaseInfo(req)
	if err != nil {
		if errors.Is(err, entity.ErrUnauthorized) {
			r.unauthorized(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	posts, comments, err := r.services.Post.GetAllCommentedPostsWithComments(userID)
	if err != nil {
		r.serverError(w, req, err)
	}

	data.Models.Posts = *posts
	for i := range data.Models.Posts {
		data.Models.Posts[i].Comments = (*comments)[i]
	}

	r.render(w, req, http.StatusOK, "commented.html", data)
}

func (r *Routes) postDelete(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	userRole := r.sesm.GetUserRole(req.Context())
	if userRole == entity.MODERATOR || userRole == entity.ADMIN {
		r.postDeletePrivileged(w, req)
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

	postID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("postDelete: invalid url path")
		r.notFound(w)
		return
	}

	err := r.services.Post.DeletePost(postID, userID)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrPostNotFound):
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

func (r *Routes) postDeletePrivileged(w http.ResponseWriter, req *http.Request) {
	userRole := r.sesm.GetUserRole(req.Context())

	postID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("postDeletePrivileged: invalid url path")
		r.notFound(w)
		return
	}

	err := r.services.Post.DeletePostPrivileged(postID, userRole)
	if err != nil {
		if errors.Is(err, entity.ErrPostNotFound) {
			r.notFound(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Routes) postReport(w http.ResponseWriter, req *http.Request) {
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

	postID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("postReport: invalid url path")
		r.notFound(w)
		return
	}
	message := req.PostFormValue("message")
	userID := r.sesm.GetUserID(req.Context())

	notification := entity.Notification{
		Type:       entity.REPORT,
		Content:    message,
		SourceID:   postID,
		SourceType: entity.POST,
		UserFrom:   userID,
	}

	err := r.services.User.SendNotification(notification)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, req.URL.Path, http.StatusOK)
}
