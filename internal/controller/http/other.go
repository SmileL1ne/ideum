package http

import (
	"fmt"
	"forum/web"
	"net/http"
)

/*
	TODO:
	1. Static handler for GET "/static/*filepath"
*/

var fileServer = http.FileServer(http.FS(web.Files))

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
	fmt.Println(posts)

	data := r.newTemplateData(req)
	data.Posts = posts

	r.render(w, req, http.StatusOK, "home.html", data)
}
