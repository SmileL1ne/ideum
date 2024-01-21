package http

import (
	"forum/internal/entity"
	"net/http"
)

func (r *routes) userSignup(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.userSignupPost(w, req)
		return
	} else if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	data := r.newTemplateData(req)
	data.Form = entity.UserSignupForm{}
	r.render(w, req, http.StatusOK, "signup.html", data)
}

func (r *routes) userSignupPost(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}

	form := req.PostForm

	username := form.Get("username")
	email := form.Get("email")
	password := form.Get("password")

	u := entity.UserSignupForm{Username: username, Email: email, Password: password}

	_, status, err := r.service.User.SaveUser(&u)
	if status != http.StatusOK {
		if status == http.StatusUnprocessableEntity {
			data := r.newTemplateData(req)
			data.Form = u
			r.render(w, req, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			r.serverError(w, req, err)
		}

		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *routes) userLogin(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Display login page"))
}

func (r *routes) userLoginPost(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Login user"))
}

func (r *routes) userLogout(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Logout user"))
}