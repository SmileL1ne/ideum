package handlers

import (
	"errors"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"strings"
)

func (r *Routes) userSignupPost(w http.ResponseWriter, req *http.Request) {
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
	email := strings.ToLower(form.Get("email"))
	password := form.Get("password")

	u := entity.UserSignupForm{Username: username, Email: email, Password: password}

	_, err := r.service.User.SaveUser(&u) // Put user id in context
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData):
			r.logger.Print("userSignupPost: invalid form fill")

			w.WriteHeader(http.StatusBadRequest)
			msg := getErrorMessage(&u.Validator)
			fmt.Fprint(w, strings.TrimSpace(msg))
		default:
			r.serverError(w, req, err)
		}
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Routes) userLoginPost(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	if err := req.ParseForm(); err != nil {
		r.logger.Print("userLoginPost: invalid form fill (parse error)")
		r.badRequest(w)
		return
	}

	form := req.PostForm
	identifier := strings.ToLower(form.Get("identifier"))
	password := form.Get("password")

	u := entity.UserLoginForm{Identifier: identifier, Password: password}
	id, err := r.service.User.Authenticate(&u)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData), errors.Is(err, entity.ErrInvalidCredentials):
			r.logger.Print("userSignupPost: invalid form fill")

			w.WriteHeader(http.StatusBadRequest)
			msg := getErrorMessage(&u.Validator)
			fmt.Fprint(w, strings.TrimSpace(msg))
		default:
			r.serverError(w, req, err)
		}
		return
	}

	// Renew session token whenever user logs in
	err = r.sesm.RenewToken(req.Context(), id)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	// Add authenticated user's id to the session data
	r.sesm.PutUserID(req.Context(), id)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func (r *Routes) userLogout(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	err := r.sesm.DeleteToken(req.Context())
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
