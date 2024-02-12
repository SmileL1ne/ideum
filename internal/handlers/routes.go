package handlers

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
*/

type routes struct {
	service   *service.Services
	tempCache map[string]*template.Template
	sesm      *sesm.SessionManager
}

// NewRouter returns http.Handler type router with registered routes
func NewRouter(s *service.Services, sesm *sesm.SessionManager) http.Handler {
	router := http.NewServeMux()

	// Temporary cache for one-time template initialization and subsequent
	// storage in templates map
	tempCache, err := newTemplateCache()
	if err != nil {
		log.Fatalf("Error creating cached templates:%v", err)
	}

	r := &routes{
		service:   s,
		tempCache: tempCache,
		sesm:      sesm,
	}

	// Serve static files
	router.HandleFunc("/static/", fileServer.ServeHTTP)

	// Dynamic middleware chain for routes that don't require authentication
	dynamic := mids.New(sesm.LoadAndSave)

	router.Handle("/", dynamic.ThenFunc(r.home))
	router.Handle("/sortByTags/", dynamic.ThenFunc(r.sortedByTag))
	router.Handle("/post/view/", dynamic.ThenFunc(r.postView)) // postID at the end
	router.Handle("/user/signup", dynamic.ThenFunc(r.userSignup))
	router.Handle("/user/login", dynamic.ThenFunc(r.userLogin))

	// Protected appends dynamic middleware chain and used for routes
	// that require authentication
	protected := dynamic.Append(r.requireAuthentication)

	router.Handle("/post/create", protected.ThenFunc(r.postCreate))
	router.Handle("/post/like/", protected.ThenFunc(r.postLike))                  // postID at the end
	router.Handle("/post/dislike/", protected.ThenFunc(r.postDislike))            // postID at the end
	router.Handle("/post/comment/", protected.ThenFunc(r.commentCreate))          // postID at the end
	router.Handle("/post/comment/like/", protected.ThenFunc(r.commentLike))       // postID at the end
	router.Handle("/post/comment/dislike/", protected.ThenFunc(r.commentDislike)) // postID at the end
	router.Handle("/user/logout", protected.ThenFunc(r.userLogout))

	// Standard middleware chain applied to router itself -> used in all routes
	standard := mids.New(r.recoverPanic, secureHeaders)

	return standard.Then(router)
}
