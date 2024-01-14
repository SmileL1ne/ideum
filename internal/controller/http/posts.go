package http

import (
	"forum/internal/entity"
	"net/http"
)

/*
	TODO:
	- Finish home handler
*/

func (r *routes) home(w http.ResponseWriter, req *http.Request) {
	info := entity.Post{}

	err := r.s.Post.SavePost(info)
	if err != nil {
		panic(err)
	}

	data := r.newTemplateData(req)

	r.render(w, req, http.StatusOK, "home.html", data)
}
