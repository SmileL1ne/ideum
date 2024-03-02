package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	RedirectURI  string
	ClientID     string
	ClientSecret string
)

func init() {

	file, err := os.Open(".env")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	variables := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			variables[parts[0]] = parts[1]
		}
	}

	RedirectURI = variables["Redirect_URI"]
	ClientID = variables["Client_id"]
	ClientSecret = variables["Client_secret"]
}

func (r *routes) googlelogin(w http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile", ClientID, RedirectURI)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

func (r *routes) googleCallback(w http.ResponseWriter, req *http.Request) {
}
