package handlers

import (
	"net/http"
	"strings"
)

// var fileServer = http.FileServer(http.FS(web.Files))

// prevenetDirListing is a middleware that prevents access to directories
// in static handler, so only full path to static files would be available
func (r *Routes) preventDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/") || len(req.URL.Path) == 0 {
			r.notFound(w)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func (r *Routes) home(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		r.notFound(w)
		return
	}

	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	data, err := r.newTemplateData(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	posts, err := r.services.Post.GetAllPosts()
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	data.Models.Posts = *posts

	r.render(w, req, http.StatusOK, "home.html", data)
}

func (r *Routes) sortedByTag(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	data, err := r.newTemplateData(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	tagID, ok := getIdFromPath(req, 3)
	if !ok {
		r.logger.Print("sortedByTag: invalid url path")
		r.notFound(w)
		return
	}

	isTagExists, err := r.services.Tag.IsExist(tagID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	if !isTagExists {
		r.logger.Print("sortedByTag: invalid tag id (doen't exist)")
		r.notFound(w)
		return
	}

	posts, err := r.services.Post.GetAllPostsByTagId(tagID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data.Models.Posts = *posts

	r.render(w, req, http.StatusOK, "home.html", data)
}
