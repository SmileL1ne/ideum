package http

import (
	"forum/web"
	"net/http"
	"strings"
)

var fileServer = http.FileServer(http.FS(web.Files))

// prevenetDirListing is a middleware that prevents access to directories
// in static handler, so only full path to static files would be available
func (r *routes) preventDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/") || len(req.URL.Path) == 0 {
			r.notFound(w)
			return
		}
		next.ServeHTTP(w, req)
	})
}

func (r *routes) home(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		r.notFound(w)
		return
	}

	if req.Method != http.MethodGet {
		r.notFound(w)
		return
	}

	posts, err := r.service.Post.GetAllPosts()
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data := r.newTemplateData(req)
	data.Models.Posts = *posts

	r.render(w, req, http.StatusOK, "home.html", data)
}
