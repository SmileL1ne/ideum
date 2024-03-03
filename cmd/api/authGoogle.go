package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/entity"
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
	code := req.URL.Query().Get("code")
	tokenURL := "https://accounts.google.com/o/oauth2/token"
	clientData := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code", code, ClientID, ClientSecret, RedirectURI)

	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", strings.NewReader(clientData))
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		r.serverError(w, req, err)
		return
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		http.Error(w, "Failed to decode token response", http.StatusInternalServerError)
		r.serverError(w, req, err)
		return
	}
	accessToken := tokenResponse["access_token"].(string)
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo"
	newreq, _ := http.NewRequest("GET", userInfoURL, nil)
	newreq.Header.Add("Authorization", "Bearer "+accessToken)

	userInfoResp, err := http.DefaultClient.Do(newreq)
	if err != nil {
		http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
		r.serverError(w, req, err)
		return
	}

	defer userInfoResp.Body.Close()

	// contents, err := ioutil.ReadAll(userInfoResp.Body)
	// if err != nil {
	// 	fmt.Fprintf(w, "failed read response: %s", err.Error())
	// 	return
	// }
	// fmt.Fprintln(w, string(contents))

	form := entity.UserSignupForm{}

	if err := json.NewDecoder(userInfoResp.Body).Decode(&form); err != nil {
		http.Error(w, "Failed to decode user info response", http.StatusInternalServerError)
		r.serverError(w, req, err)
		return
	}

	form.Password, err = RandomPassword(8)
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		r.serverError(w, req, err)
		return
	}

	fmt.Println("inFO", form)


	r.Sso(w, req, form)
}

func (r *routes) Sso(w http.ResponseWriter, req *http.Request, form entity.UserSignupForm) {
	id, err := r.service.User.SaveUser(&form)

	fmt.Println("OPA", form)
	if err != nil {
		fmt.Println("OI", err)
		switch {
		case errors.Is(err, entity.ErrDuplicateEmail) || errors.Is(err, entity.ErrDuplicateUsername):

			user, err := r.service.User.GetUserByEmail(form.Email)
			fmt.Println("-", user.Id)
			if err != nil {
				fmt.Println("line109")
				r.serverError(w, req, err)
				return
			}

			err = r.sesm.RenewToken(req.Context(), user.Id)
			if err != nil {
				fmt.Println("line116")
				r.serverError(w, req, err)
				return
			}

			r.sesm.PutUserID(req.Context(), user.Id)
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return

		default:
			r.serverError(w, req, err)
			fmt.Println(err)
			fmt.Println("line 127")
			return
		}

	}
	// login here

	// Renew session token whenever user logs in
	err = r.sesm.RenewToken(req.Context(), id)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	// Add authenticated user's id to the session data
	r.sesm.PutUserID(req.Context(), id)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

func RandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}