package http

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
)

/*
	TODO:
	- Create error page func that would nicely render error pages (refer - 'ERR')
*/

func (r *routes) newTemplateData(req *http.Request) templateData {
	return templateData{
		Flash:           r.sesm.PopString(req.Context(), "flash"),
		IsAuthenticated: r.isAuthenticated(req),
	}
}

func (r *routes) isAuthenticated(req *http.Request) bool {
	return r.sesm.Exists(req.Context(), "authenticatedUserID")
}

func (r *routes) serverError(w http.ResponseWriter, req *http.Request, err error) {
	var (
		method = req.Method
		uri    = req.RequestURI
		trace  = string(debug.Stack())
	)

	log.Printf(err.Error()+" method - %s, uri - %s, stack - %s", method, uri, trace)
	/*
		ERR
	*/
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (r *routes) clientError(w http.ResponseWriter, status int) {
	/*
		ERR
	*/
	http.Error(w, http.StatusText(status), status)
}

func (r *routes) notFound(w http.ResponseWriter) {
	r.clientError(w, http.StatusNotFound)
}

func (r *routes) badRequest(w http.ResponseWriter) {
	r.clientError(w, http.StatusBadRequest)
}

func (r *routes) methodNotAllowed(w http.ResponseWriter) {
	r.clientError(w, http.StatusMethodNotAllowed)
}

// Render templates by retrieving necessary template from template cache.
//
// First execute into dummy buffer for any execution error catch (to set appropriate header)
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

// getIdFromPath retrieves and returns id from request path.
//
// It returns empty string if number of splitted parts doesn't match with
// given number
func (r *routes) getIdFromPath(req *http.Request, urlPartsNum int) string {
	urlParts := strings.Split(req.URL.Path, "/")
	if len(urlParts) != urlPartsNum {
		return ""
	}

	return urlParts[len(urlParts)-1]
}
