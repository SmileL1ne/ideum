package handlers

import (
	"errors"
	"fmt"
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
		r.serverError(w, req, err)
		return
	}

	requests, err := r.services.User.GetRequests()
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data.Models.Requests = *requests

	r.render(w, req, http.StatusOK, "request.html", data)
}

func (r *Routes) reports(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		r.methodNotAllowed(w)
		return
	}
	data, err := r.newTemplateData(req)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	reports, err := r.services.User.GetReports()
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	data.Models.Reports = *reports

	r.render(w, req, http.StatusOK, "report.html", data)
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

	userTo, ok := getIdFromPath(req, 4)
	if !ok {
		r.logger.Print("promoteUser: invalid url path")
		r.notFound(w)
		return
	}

	promotionID, ok := getValidID(req.PostForm.Get("promotionID"))
	if !ok {
		r.logger.Print("promoteUser: invalid promotionID")
		r.badRequest(w)
		return
	}

	err := r.services.User.DeletePromotion(promotionID)
	if err != nil {
		if errors.Is(err, entity.ErrPromotionNotFound) {
			r.notFound(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	err = r.services.User.PromoteUser(userTo)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	adminID := r.sesm.GetUserID(req.Context())

	notification := entity.Notification{
		Type:     entity.PROMOTED,
		UserFrom: adminID,
		UserTo:   userTo,
	}

	err = r.services.User.SendNotification(notification)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDuplicateNotification):
			r.logger.Print("promoteUser: duplicate notification")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Promotion of user notification is already sent")
		default:
			r.serverError(w, req, err)
		}
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

	promotionID, ok := getValidID(req.PostForm.Get("promotionID"))
	if !ok {
		r.logger.Print("rejectPromotion: invalid promotionID")
		r.badRequest(w)
		return
	}

	err := r.services.User.DeletePromotion(promotionID)
	if err != nil {
		if errors.Is(err, entity.ErrPromotionNotFound) {
			r.notFound(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	adminID := r.sesm.GetUserID(req.Context())

	notification := entity.Notification{
		Type:     entity.REJECT_PROMOTION,
		UserFrom: adminID,
		UserTo:   userID,
	}

	err = r.services.User.SendNotification(notification)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDuplicateNotification):
			r.logger.Print("rejectPromotion: duplicate notification")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Reject promotion notification is already sent")
		default:
			r.serverError(w, req, err)
		}
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

	reportID, ok := getValidID(req.PostForm.Get("reportID"))
	if !ok {
		r.logger.Print("rejectReport: invalid reportID")
		r.badRequest(w)
		return
	}

	err := r.services.User.DeleteReport(reportID)
	if err != nil {
		if errors.Is(err, entity.ErrReportNotFound) {
			r.notFound(w)
			return
		}
		r.serverError(w, req, err)
		return
	}

	sourceID, err := strconv.Atoi(req.PostForm.Get("sourceID"))
	if err != nil {
		r.logger.Print("rejectReport:", err)
		r.badRequest(w)
		return
	}

	adminID := r.sesm.GetUserID(req.Context())

	notification := entity.Notification{
		Type:     entity.REJECT_REPORT,
		SourceID: sourceID,
		UserFrom: adminID,
		UserTo:   userID,
	}

	err = r.services.User.SendNotification(notification)
	if err != nil {
		switch {
		case errors.Is(err, entity.ErrDuplicateNotification):
			r.logger.Print("rejectReport: duplicate notification")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, "Reject report notification is already sent")
		default:
			r.serverError(w, req, err)
		}
		return
	}

	http.Redirect(w, req, "/user/notifications", http.StatusSeeOther)
}
