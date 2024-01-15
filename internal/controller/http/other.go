package http

import (
	"forum/internal/entity"
	"forum/web"
	"net/http"
)

/*
	TODO:
	1. Static handler for GET "/static/*filepath"
*/

var fileServer = http.FileServer(http.FS(web.Files))

func (r *routes) home(w http.ResponseWriter, req *http.Request) {
	info := entity.Post{}

	err := r.s.Post.SavePost(info)
	if err != nil {
		panic(err)
	}

	data := r.newTemplateData(req)

	r.render(w, req, http.StatusOK, "home.html", data)
}
