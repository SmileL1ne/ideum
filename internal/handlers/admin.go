package handlers

import (
	"errors"
	"forum/internal/entity"
	"net/http"
)

func (r *Routes) requests(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}
	data, err := r.newTemplateData(req)
	if err != nil {
		if errors.Is(err, entity.ErrUnauthorized) {
			r.unauthorized(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	userRole := r.sesm.GetUserRole(req.Context())
	if userRole != entity.ADMIN {
		r.forbidden(w)
		return
	}

	requests, err := r.services.User.GetRequests(userRole)
	if err != nil {
		if errors.Is(err, entity.ErrForbiddenAccess) {
			r.forbidden(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	data.Models.Notifications = *requests

	r.render(w, req, http.StatusOK, "notification.html", data)
}

func (r *Routes) adminPromote(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("adminPromote: invalid url path")
		r.notFound(w)
		return
	}

	err := r.services.User.PromoteUser(userID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	notification := entity.Notification{
		Type:   entity.PROMOTED,
		UserTo: userID,
	}

	err = r.services.User.SendNotification(notification)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, "/admin/requests", http.StatusSeeOther)
}

func (r *Routes) adminReject(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}

	userID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("adminPromote: invalid url path")
		r.notFound(w)
		return
	}

	notification := entity.Notification{
		Type:   entity.REJECT_PROMOTION,
		UserTo: userID,
	}

	err := r.services.User.SendNotification(notification)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, "/admin/requests", http.StatusSeeOther)
}
