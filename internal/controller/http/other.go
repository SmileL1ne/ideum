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

	data := r.newTemplateData(req)

	r.render(w, req, http.StatusOK, "home.html", data)
}
