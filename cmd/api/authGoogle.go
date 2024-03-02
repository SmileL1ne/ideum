package main

import (
	"fmt"
	"net/http"
)

const (
	Redirect_URI  = "https://localhost:5000/callbackGoogle"
	Client_id     = "475115590650-7kmdkikv6tfhh0kfia3s2hcvfpffi5re.apps.googleusercontent.com"
	Client_secret = "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"
)

func (r *routes) googlelogin(w http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile", Client_id, Redirect_URI)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

func (r *routes) googleCallback(w http.ResponseWriter, req *http.Request) {
}
