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

	_, err := r.services.User.SaveUser(&u) // Put user id in context
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrInvalidFormData), errors.Is(err, entity.ErrDuplicateUsername), errors.Is(err, entity.ErrDuplicateEmail):
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
	id, err := r.services.User.Authenticate(&u)
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

	role, err := r.services.User.GetUserRole(id)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	err = r.sesm.RenewToken(req.Context(), id)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	r.sesm.PutUserID(req.Context(), id)
	r.sesm.PutUserRole(req.Context(), role)

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

func (r *Routes) notifications(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID, data, err := r.getBaseInfo(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	notifications, err := r.services.User.GetNotifications(userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data.Models.Notifications = *notifications

	r.render(w, req, http.StatusOK, "notification.html", data)
}

func (r *Routes) deleteNotification(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}

	notificationID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("deleteNotification: invalid url path")
		r.notFound(w)
		return
	}

	err := r.services.User.DeleteNotification(notificationID)
	if err != nil {
		if errors.Is(err, entity.ErrNotificationNotFound) {
			r.notFound(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, "/user/notifications", http.StatusSeeOther)
}

func (r *Routes) userPromote(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}
	if r.sesm.GetUserRole(req.Context()) != entity.USER {
		r.badRequest(w)
		return
	}

	userID := r.sesm.GetUserID(req.Context())

	err := r.services.User.SendPromotion(userID)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicatePromotion) {
			r.logger.Print("userPromote:", err)
			r.badRequest(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, "/", http.StatusSeeOther)
}
