package http

import (
	"forum/internal/service"
	"forum/pkg/mids"
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

type chain struct {
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

	dynamic := mids.New(sesm.LoadAndSave)

	router.Handle("/", dynamic.Then(http.HandlerFunc(r.home)))
	router.Handle("/post/view/", dynamic.Then(http.HandlerFunc(r.postView)))
	router.Handle("/user/signup", dynamic.Then(http.HandlerFunc(r.userSignup)))
	router.Handle("/user/signup/post", dynamic.Then(http.HandlerFunc(r.userSignupPost)))
	router.Handle("/user/login", dynamic.Then(http.HandlerFunc(r.userLogin)))

	protected := dynamic.Append(r.requireAuthentication)

	router.Handle("/post/create", protected.Then(http.HandlerFunc(r.postCreate)))
	router.Handle("/post/create/post", protected.Then(http.HandlerFunc(r.postCreatePost)))
	router.Handle("/user/logout", protected.Then(http.HandlerFunc(r.userLogout)))

	standard := mids.New(secureHeaders)

	return standard.Then(router)
}
