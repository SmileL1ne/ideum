package main

import (
	"bytes"
	"errors"
	"fmt"
	"forum/internal/entity"
	"log"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

/*
	TODO:
	- Create error page func that would render prettys error pages (refer to 'ERR')
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

	log.Printf(err.Error()+"; method - %s, uri - %s, stack - %s", method, uri, trace)
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

func (r *routes) unauthorized(w http.ResponseWriter) {
	r.clientError(w, http.StatusUnauthorized)
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

// getUsername retrieves username by user id from request context
func (r *routes) getUsername(req *http.Request) (string, error) {
	userID := r.sesm.GetInt(req.Context(), "authenticatedUserID")
	if userID == 0 {
		return "", entity.ErrInvalidUserID
	}

	username, err := r.service.User.GetUsernameByID(userID)
	if errors.Is(err, entity.ErrInvalidCredentials) {
		return "", entity.ErrInvalidUserID
	}
	return username, err
}

// getIdFromPath retrieves and returns id from request path.
//
// It returns empty string if number of splitted parts doesn't match with
// given number
func getIdFromPath(req *http.Request, urlPartsNum int) (int, error) {
	urlParts := strings.Split(req.URL.Path, "/")
	if len(urlParts) != urlPartsNum {
		return 0, entity.ErrInvalidURLPath
	}

	id, isValid := getValidID(urlParts[len(urlParts)-1])
	if !isValid {
		return 0, entity.ErrInvalidPathID
	}

	return id, nil
}

// getValidID parses string id to int and checks if it is valid
func getValidID(idStr string) (int, bool) {
	id, _ := strconv.Atoi(idStr)
	if id <= 0 {
		return 0, false
	}
	return id, true
}
