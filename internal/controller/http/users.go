package http

import (
	"errors"
	"forum/internal/entity"
	"net/http"
)

func (r *routes) userSignup(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	data := r.newTemplateData(req)
	data.Form = entity.UserSignupForm{}
	r.render(w, req, http.StatusOK, "signup.html", data)
}

func (r *routes) userSignupPost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}

	form := req.PostForm

	username := form.Get("username")
	email := form.Get("email")
	password := form.Get("password")

	u := entity.UserSignupForm{Username: username, Email: email, Password: password}

	_, err := r.service.User.SaveUser(&u) // Put user id in context
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			data := r.newTemplateData(req)
			data.Form = u
			r.render(w, req, http.StatusUnprocessableEntity, "signup.html", data)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	r.sesm.Put(req.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, req, "/user/login", http.StatusSeeOther)
}

func (r *routes) userLogin(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		r.userLoginPost(w, req)
		return
	} else if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	data := r.newTemplateData(req)
	data.Form = entity.UserLoginForm{}
	r.render(w, req, http.StatusOK, "login.html", data)
}

func (r *routes) userLoginPost(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}

	form := req.PostForm
	identifier := form.Get("identifier")
	password := form.Get("password")

	u := entity.UserLoginForm{Identifier: identifier, Password: password}
	id, err := r.service.User.Authenticate(&u)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidCredentials):
			u.AddNonFieldError("Email or password is incorrect")
			data := r.newTemplateData(req)
			data.Form = u
			r.render(w, req, http.StatusUnprocessableEntity, "login.html", data)
		default:
			r.serverError(w, req, err)
		}
		return
	}

	err = r.sesm.RenewToken(req.Context())
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	r.sesm.Put(req.Context(), "authenticatedUserID", id)
	r.sesm.Put(req.Context(), "flash", "Successfully logged in!")

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *routes) userLogout(w http.ResponseWriter, req *http.Request) {
	err := r.sesm.RenewToken(req.Context())
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	r.sesm.Remove(req.Context(), "authenticatedUserID")
	r.sesm.Put(req.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
