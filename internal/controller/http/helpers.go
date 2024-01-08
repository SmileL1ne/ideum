package http

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// TODO: Create error page func that would nicely render error pages (refer - 'ERR')

func (r *routes) newTemplateData(req *http.Request) templateData {
	return templateData{}
}

func (r *routes) serverError(w http.ResponseWriter, req *http.Request, err error) {
	var (
		method = req.Method
		uri    = req.RequestURI
		trace  = string(debug.Stack())
	)

	r.l.Error(err.Error(), "method", method, "uri", uri, "stack", trace)
	// ERR
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (r *routes) clientError(w http.ResponseWriter, status int) {
	// ERR
	http.Error(w, http.StatusText(status), status)
}

func (r *routes) render(w http.ResponseWriter, req *http.Request, status int, page string, data templateData) {
	tmpl, ok := r.tempCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		r.serverError(w, req, err)
		return
	}

	buf := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(buf, "base", data); err != nil {
		r.serverError(w, req, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}
