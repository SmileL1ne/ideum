package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"forum/internal/entity"
	"forum/internal/validator"
	"forum/web"
	"html/template"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

func (r *Routes) newTemplateData(req *http.Request) templateData {
	return templateData{
		IsAuthenticated: r.isAuthenticated(req),
	}
}

func (r *Routes) isAuthenticated(req *http.Request) bool {
	return r.sesm.ExistsUserID(req.Context())
}

func (r *Routes) serverError(w http.ResponseWriter, req *http.Request, err error) {
	var (
		method = req.Method
		uri    = req.RequestURI
		trace  = string(debug.Stack())
	)

	r.logger.Printf(err.Error()+"; method - %s, uri - %s, stack - %s", method, uri, trace)

	errInfo := errData{
		ErrCode: http.StatusInternalServerError,
		ErrMsg:  http.StatusText(http.StatusInternalServerError),
	}

	// Custom render of error template for server error
	// (to avoid infinite recursion of original render function)
	r.renderErrorPage(w, errInfo)
}

func (r *Routes) clientError(w http.ResponseWriter, status int) {
	errInfo := errData{
		ErrCode: status,
		ErrMsg:  http.StatusText(status),
	}

	r.renderErrorPage(w, errInfo)
}

func (r *Routes) unauthorized(w http.ResponseWriter) {
	r.clientError(w, http.StatusUnauthorized)
}

func (r *Routes) notFound(w http.ResponseWriter) {
	r.clientError(w, http.StatusNotFound)
}

func (r *Routes) badRequest(w http.ResponseWriter) {
	r.clientError(w, http.StatusBadRequest)
}

func (r *Routes) methodNotAllowed(w http.ResponseWriter) {
	r.clientError(w, http.StatusMethodNotAllowed)
}

// Render templates by retrieving necessary template from template cache.
//
// First execute into dummy buffer for any execution error catch (to set appropriate header)
func (r *Routes) render(w http.ResponseWriter, req *http.Request, status int, page string, data templateData) {
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

// renderErrorPage renders error page with given error code and message
func (r *Routes) renderErrorPage(w http.ResponseWriter, errInfo errData) {
	tmpl, err := template.ParseFS(web.Files, "html/error.html")
	if err != nil {
		r.logger.Print("the template error.html does not exist")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, errInfo); err != nil {
		r.logger.Print("error executing error.html template")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(errInfo.ErrCode)

	buf.WriteTo(w)
}

// getBaseInfo retrieves all information for base of the page (username and all tags)
func (r *Routes) getBaseInfo(req *http.Request) (string, *[]entity.TagEntity, error) {
	username, err := r.getUsername(req)
	if err != nil {
		return "", nil, err
	}

	tags, err := r.service.Tag.GetAllTags()
	if err != nil {
		return "", nil, err
	}

	return username, tags, nil
}

// getUsername retrieves username by user id from request context
func (r *Routes) getUsername(req *http.Request) (string, error) {
	userID := r.sesm.GetUserID(req.Context())
	if userID == 0 {
		return "", nil
	}

	username, err := r.service.User.GetUsernameById(userID)
	if errors.Is(err, entity.ErrInvalidCredentials) {
		return "", nil
	}
	return username, err
}

// getIdFromPath retrieves and returns id from request path.
//
// It returns empty string if number of splitted parts doesn't match with
// given number
func getIdFromPath(req *http.Request, urlPartsNum int) (int, bool) {
	urlParts := strings.Split(req.URL.Path, "/")
	if len(urlParts) != urlPartsNum {
		return 0, false
	}

	id, isValid := getValidID(urlParts[len(urlParts)-1])
	if !isValid {
		return 0, false
	}

	return id, true
}

// getValidID parses string id to int and checks if it is valid
func getValidID(idStr string) (int, bool) {
	id, err := strconv.Atoi(idStr)
	if id < 0 || err != nil {
		return 0, false
	}
	return id, true
}

// getErrorMessage accepts pointer to form's validator that should consist of
// field and/or non field errors and returns formatted error message
func getErrorMessage(v *validator.Validator) string {
	var msg string

	for _, str := range v.NonFieldErrors {
		msg += str + "\n"
	}
	for key, val := range v.FieldErrors {
		msg += key + ": " + val + "\n"
	}

	return msg
}
