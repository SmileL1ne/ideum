package http

import "net/http"

func (r *routes) userSignup(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Display user signup page"))
}

func (r *routes) userSignupPost(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Signup a user (create a new user)"))
}

func (r *routes) userLogin(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Display login page"))
}

func (r *routes) userLoginPost(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Login user"))
}

func (r *routes) userLogout(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Logout user"))
}
