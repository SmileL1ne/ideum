package handlers

import (
	"errors"
	"forum/internal/entity"
	"net/http"
	"strconv"
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

func (r *Routes) promoteUser(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}

	userID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("promoteUser: invalid url path")
		r.notFound(w)
		return
	}

	notificationID, ok := getValidID(req.PostForm.Get("notificationID"))
	if !ok {
		r.logger.Print("promoteUser: invalid notificationID")
		r.badRequest(w)
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

	err = r.services.User.DeleteNotification(notificationID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, "/admin/requests", http.StatusSeeOther)
}

func (r *Routes) rejectPromotion(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}

	userID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("rejectPromotion: invalid url path")
		r.notFound(w)
		return
	}

	notificationID, ok := getValidID(req.PostForm.Get("notificationID"))
	if !ok {
		r.logger.Print("promoteUser: invalid notificationID")
		r.badRequest(w)
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

	err = r.services.User.DeleteNotification(notificationID)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	http.Redirect(w, req, "/admin/requests", http.StatusSeeOther)
}

func (r *Routes) rejectReport(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		r.methodNotAllowed(w)
		return
	}
	if err := req.ParseForm(); err != nil {
		r.badRequest(w)
		return
	}

	userID, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("rejectReport: invalid url path")
		r.notFound(w)
		return
	}

	notificationID, ok := getValidID(req.PostForm.Get("notificationID"))
	if !ok {
		r.logger.Print("promoteUser: invalid notificationID")
		r.badRequest(w)
		return
	}

	sourceID, err := strconv.Atoi(req.PostForm.Get("sourceID"))
	if err != nil {
		r.logger.Print("rejectReport:", err)
		r.badRequest(w)
		return
	}

	sourceType := req.PostForm.Get("sourceType")
	if sourceType == "" || (sourceType != entity.POST && sourceType != entity.COMMENT) {
		r.logger.Print("rejectReport: invalid sourceType provided")
		r.badRequest(w)
		return
	}

	notification := entity.Notification{
		Type:       entity.REJECT_REPORT,
		SourceType: sourceType,
		SourceID:   sourceID,
		UserTo:     userID,
	}

	err = r.services.User.SendNotification(notification)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	err = r.services.User.DeleteNotification(notificationID)
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
