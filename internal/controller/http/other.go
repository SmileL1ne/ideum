package http

import (
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
		r.clientError(w, http.StatusNotFound)
		return
	}
	if req.Method != http.MethodGet {
		r.clientError(w, http.StatusNotFound)
		return
	}

	data := r.newTemplateData(req)

	r.render(w, req, http.StatusOK, "home.html", data)
}
