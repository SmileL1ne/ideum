package handlers

import (
	"forum/internal/service"
	"forum/pkg/mids"
	"forum/pkg/sesm"
	"html/template"
	"log"
	"net/http"
)

type Routes struct {
	service   *service.Services
	tempCache map[string]*template.Template
	sesm      *sesm.SessionManager
	logger    *log.Logger
}

func NewRouter(
	services *service.Services,
	sesm *sesm.SessionManager,
	logger *log.Logger,
) *Routes {

	// Temporary cache for one-time template initialization and subsequent
	// storage in templates map
	tempCache, err := newTemplateCache()
	if err != nil {
		log.Fatalf("Error creating cached templates:%v", err)
	}

	return &Routes{
		service:   services,
		tempCache: tempCache,
		sesm:      sesm,
		logger:    logger,
	}
}

// NewRouter returns http.Handler type router with registered routes
func (r *Routes) Register() http.Handler {
	router := http.NewServeMux()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./web/static/"))
	router.Handle("/static/", r.preventDirListing(http.StripPrefix("/static", fileServer)))

	// Dynamic middleware chain for routes that don't require authentication
	dynamic := mids.New(r.sesm.LoadAndSave)

	router.Handle("/", dynamic.ThenFunc(r.home))
	router.Handle("/sortByTags/", dynamic.ThenFunc(r.sortedByTag))
	router.Handle("/user/login", dynamic.ThenFunc(r.userLoginPost))
	router.Handle("/user/signup", dynamic.ThenFunc(r.userSignupPost))
	router.Handle("/post/view/", dynamic.ThenFunc(r.postView)) // postID at the end

	// Protected appends dynamic middleware chain and used for routes
	// that require authentication
	protected := dynamic.Append(r.requireAuthentication)

	router.Handle("/post/myPosts", protected.ThenFunc(r.postsPersonal))
	router.Handle("/post/myReacted", protected.ThenFunc(r.postsReacted))
	router.Handle("/post/create", protected.ThenFunc(r.postCreate))
	router.Handle("/post/reaction/", protected.ThenFunc(r.postReaction))            // postID at the end
	router.Handle("/post/comment/", protected.ThenFunc(r.commentCreate))            // postID at the end
	router.Handle("/post/comment/reaction/", protected.ThenFunc(r.commentReaction)) // postID at the end
	router.Handle("/user/logout", protected.ThenFunc(r.userLogout))

	// Standard middleware chain applied to router itself -> used in all routes
	standard := mids.New(r.recoverPanic, secureHeaders)

	return standard.Then(router)
}
