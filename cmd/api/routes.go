package main

import (
	"forum/pkg/mids"
	"net/http"
)

// NewRouter returns http.Handler type router with registered routes
func (r *routes) newRouter() http.Handler {
	router := http.NewServeMux()

	// Serve static files
	router.Handle("/static/", r.preventDirListing(fileServer))

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
