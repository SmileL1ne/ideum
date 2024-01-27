package http

import (
	"forum/internal/service"
	"forum/pkg/sesm"
	"html/template"
	"log"
	"net/http"
)

/*
	TODO:
	- Checkout comments in each registered handlers
	-- Add middlewares
*/

type routes struct {
	service   *service.Service
	tempCache map[string]*template.Template
	sesm      *sesm.SessionManager
}

func NewRouter(s *service.Service, sesm *sesm.SessionManager) http.Handler {
	router := http.NewServeMux()

	tempCache, err := newTemplateCache()
	if err != nil {
		log.Fatalf("Error creating cached templates:%v", err)
	}

	r := &routes{
		service:   s,
		tempCache: tempCache,
		sesm:      sesm,
	}

	router.HandleFunc("/static/", fileServer.ServeHTTP)

	router.Handle("/", r.sesm.LoadAndSave(http.HandlerFunc(r.home)))

	router.Handle("/post/view/", r.sesm.LoadAndSave(http.HandlerFunc(r.postView)))
	router.Handle("/post/create", r.sesm.LoadAndSave(http.HandlerFunc(r.postCreate)))
	router.Handle("/post/create/post", r.sesm.LoadAndSave(http.HandlerFunc(r.postCreatePost)))

	router.Handle("/user/signup", r.sesm.LoadAndSave(http.HandlerFunc(r.userSignup)))
	router.Handle("/user/signup/post", r.sesm.LoadAndSave(http.HandlerFunc(r.userSignupPost)))
	router.Handle("/user/login", r.sesm.LoadAndSave(http.HandlerFunc(r.userLogin)))
	router.Handle("/user/logout", r.sesm.LoadAndSave(http.HandlerFunc(r.userLogout)))

	return router
}
