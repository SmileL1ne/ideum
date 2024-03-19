package handlers

import (
	"forum/config"
	"forum/internal/service"
	"forum/pkg/mids"
	"forum/pkg/sesm"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

type userRateLimit struct {
	lastReq time.Time
	penalty time.Duration
}

type Routes struct {
	services       *service.Services
	tempCache      map[string]*template.Template
	sesm           *sesm.SessionManager
	logger         *log.Logger
	cfg            *config.Config
	userRateLimits map[string]userRateLimit
	rateMu         *sync.Mutex
}

func NewRouter(
	services *service.Services,
	sesm *sesm.SessionManager,
	logger *log.Logger,
	cfg *config.Config,
) *Routes {

	// Temporary cache for one-time template initialization and subsequent
	// storage in templates map
	tempCache, err := newTemplateCache()
	if err != nil {
		log.Fatalf("Error creating cached templates:%v", err)
	}

	return &Routes{
		services:       services,
		tempCache:      tempCache,
		sesm:           sesm,
		logger:         logger,
		cfg:            cfg,
		userRateLimits: make(map[string]userRateLimit),
		rateMu:         &sync.Mutex{},
	}
}

// NewRouter returns http.Handler type router with registered routes
func (r *Routes) Register() http.Handler {
	router := http.NewServeMux()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./web/static/"))
	router.Handle("/static/", r.preventDirListing(http.StripPrefix("/static", fileServer)))

	// Dynamic middleware chain for routes that don't require authentication
	dynamic := mids.New(r.sesm.LoadAndSave, r.detectUserRole)

	router.Handle("/", dynamic.ThenFunc(r.home))
	router.Handle("/sortByTags/", dynamic.ThenFunc(r.sortedByTag))
	router.Handle("/user/login", dynamic.ThenFunc(r.userLoginPost))
	router.Handle("/user/signup", dynamic.ThenFunc(r.userSignupPost))
	router.Handle("/post/view/", dynamic.ThenFunc(r.postView)) // postID at the end

	// SSO
	router.Handle("/login/google", dynamic.ThenFunc(r.googlelogin))
	router.Handle("/callbackGoogle", dynamic.ThenFunc(r.googleCallback))
	//git hub sso
	router.Handle("/login/github", dynamic.ThenFunc(r.githublogin))
	router.Handle("/callbackGithub", dynamic.ThenFunc(r.githubCallback))

	// Protected appends dynamic middleware chain and used for routes
	// that require authentication
	protected := dynamic.Append(r.requireAuthentication)

	router.Handle("/post/myPosts", protected.ThenFunc(r.postsPersonal))
	router.Handle("/post/myReacted", protected.ThenFunc(r.postsReacted))
	router.Handle("/post/myCommented", protected.ThenFunc(r.postsCommented))
	router.Handle("/post/create", protected.ThenFunc(r.postCreate))
	router.Handle("/post/delete/", protected.ThenFunc(r.postDelete))                // postID at the end
	router.Handle("/post/report/", protected.ThenFunc(r.postReport))                // postID at the end
	router.Handle("/post/reaction/", protected.ThenFunc(r.postReaction))            // postID at the end
	router.Handle("/post/comment/", protected.ThenFunc(r.commentCreate))            // postID at the end
	router.Handle("/post/comment/reaction/", protected.ThenFunc(r.commentReaction)) // postID at the end
	router.Handle("/admin/requests", protected.ThenFunc(r.requests))
	router.Handle("/admin/promote/", protected.ThenFunc(r.adminPromote)) // userID at the end
	router.Handle("/admin/reject/", protected.ThenFunc(r.adminReject))   // userID at the end
	router.Handle("/user/promote", protected.ThenFunc(r.userPromote))
	router.Handle("/user/logout", protected.ThenFunc(r.userLogout))

	// Standard middleware chain applied to router itself -> used in all routes
	standard := mids.New(r.recoverPanic, r.limitRate, r.secureHeaders)

	return standard.Then(router)
}
