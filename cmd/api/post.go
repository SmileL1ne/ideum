package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"forum/internal/entity"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (r *routes) postView(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	username, tags, err := r.getBaseInfo(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	postID, err := getIdFromPath(req, 4)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidURLPath):
			r.logger.Print("postView: invalid url path")
			r.notFound(w)
		case errors.Is(err, entity.ErrInvalidPathID):
			r.logger.Print("postView: invalid id in request path")
			r.badRequest(w)
		}
		return
	}

	post, err := r.service.Post.GetPost(postID)
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

	comments, err := r.service.Comment.GetAllCommentsForPost(postID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Post = post
	data.Models.Comments = *comments
	data.Models.Tags = *tags
	data.Username = username

	r.render(w, req, http.StatusOK, "view.html", data)
}

func (r *routes) postCreate(w http.ResponseWriter, req *http.Request) {
	switch {
	case req.Method == http.MethodPost:
		r.postCreatePost(w, req)
		return
	case req.Method != http.MethodGet:
		r.methodNotAllowed(w)
		return
	}

	username, tags, err := r.getBaseInfo(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Tags = *tags
	data.Username = username

	r.render(w, req, http.StatusOK, "create.html", data)
}

func (r *routes) postCreatePost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	// 10 means that file will be stored temporary 10 by
	if err := req.ParseMultipartForm(10); err != nil {
		r.logger.Print("postCreatePost: invalid form fill (parse error)")
		r.badRequest(w)
		return
	}

	form := req.PostForm
	title := form.Get("title")
	content := form.Get("content")

	// Get all selected tags id
	tags := form["tags"]
	if len(tags) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		msg := "tags: At least one tag should be selected"
		fmt.Fprint(w, strings.TrimSpace(msg))
		return
	}

	// Get userID from request
	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	p := entity.PostCreateForm{Title: title, Content: content}

	id, err := r.service.Post.SavePost(&p, userID, tags)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			r.logger.Print("postCreatePost: invalid form fill")

			w.WriteHeader(http.StatusBadRequest)
			msg := getErrorMessage(&p.Validator)
			fmt.Fprint(w, strings.TrimSpace(msg))
		default:
			r.serverError(w, req, err)
		}
		return
	}

	// Save image

	// take from file
	file, fileHeader, err := req.FormFile("image")
	if err != nil {
		r.logger.Print("postCreatePost: no file")

	} else {

		defer file.Close()

		// check content extention
		contentType := fileHeader.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" && contentType != "image/jpg" {
			w.WriteHeader(http.StatusBadRequest)
			msg := "only jpeg,png,gif allowed"
			fmt.Fprint(w, strings.TrimSpace(msg))
			return
		}

		// check file size 20 << 20 |  20* 1024 * 1024
		if fileHeader.Size > (20 << 20) {
			w.WriteHeader(http.StatusBadRequest)
			msg := "file size too big"
			fmt.Fprint(w, strings.TrimSpace(msg))
			return
		}

		//create sha256 file name to store one file if users loads the same file
		ext := strings.Split(fileHeader.Filename, ".")[1]
		h := sha256.New()
		io.Copy(h, file)
		fname := fmt.Sprintf("%x", h.Sum(nil)) + "." + ext
		path := filepath.Join("./web/static/public/", fname)

		newFile, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			r.logger.Print("postCreatePost: cant create file")
			r.badRequest(w) // TODO: internal server error?
			return
		}
		defer newFile.Close()

		file.Seek(0, 0)
		io.Copy(newFile, file)
		// insert to table
		if err := r.service.Image.SaveImage(id, fname); err != nil {
			r.logger.Print("Image no added to db")
		}

	}
	
	redirectURL := fmt.Sprintf("/post/view/%d", id)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, redirectURL)
}

func (r *routes) postsPersonal(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	username, tags, err := r.getBaseInfo(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if username == "" { // This should never happen
		r.unauthorized(w)
		return
	}

	userPosts, err := r.service.Post.GetAllPostsByUserId(userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Tags = *tags
	data.Models.Posts = *userPosts
	data.Username = username

	r.render(w, req, http.StatusOK, "home.html", data)
}

func (r *routes) postsReacted(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		r.unauthorized(w)
		return
	}

	username, tags, err := r.getBaseInfo(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if username == "" {
		r.unauthorized(w)
		return
	}

	reactedPosts, err := r.service.Post.GetAllPostsByUserReaction(userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Tags = *tags
	data.Models.Posts = *reactedPosts
	data.Username = username

	r.render(w, req, http.StatusOK, "home.html", data)
}
