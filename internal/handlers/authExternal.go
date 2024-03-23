package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/entity"
	"forum/pkg/pswd"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	ACCESS_TOKEN         = "access_token"
	GOOGLE_USER_INFO_URL = "https://www.googleapis.com/oauth2/v2/userinfo"
	GOOGLE_TOKEN_URL     = "https://accounts.google.com/o/oauth2/token"
	GITHUB_USER_INFO_URL = "https://api.github.com/user"
	GITHUB_TOKEN_URL     = "https://github.com/login/oauth/access_token"
)

// SSO is external auth login handler that handles registration (or authorization if user
// already exists) to the website
func (r *Routes) SSO(w http.ResponseWriter, req *http.Request, form *entity.UserSignupForm) {
	id, err := r.services.User.SaveUser(form)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicateEmail) || errors.Is(err, entity.ErrDuplicateUsername) {
			user, err := r.services.User.GetUserByEmail(form.Email)
			if err != nil {
				r.serverError(w, req, err)
				return
			}

			id = user.Id
		} else {
			r.serverError(w, req, err)
			return
		}
	}

	role, err := r.services.User.GetUserRole(id)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	err = r.sesm.RenewToken(req.Context(), id)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	r.sesm.PutUserID(req.Context(), id)
	r.sesm.PutUserRole(req.Context(), role)

	http.Redirect(w, req, "/", http.StatusSeeOther)
}

/*
	GOOGLE LOGIN
*/

func (r *Routes) googlelogin(w http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile",
		r.cfg.ExternalAuth.GoogleClientID, r.cfg.ExternalAuth.GoogleRedirectURL)

	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

func (r *Routes) googleCallback(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	clientData := fmt.Sprintf("code=%s&client_id=%s&client_secret=%s&redirect_uri=%s&grant_type=authorization_code",
		code, r.cfg.ExternalAuth.GoogleClientID, r.cfg.ExternalAuth.GoogleClientSecret, r.cfg.ExternalAuth.GoogleRedirectURL)

	resp, err := http.Post(GOOGLE_TOKEN_URL, "application/x-www-form-urlencoded", strings.NewReader(clientData))
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		r.serverError(w, req, err)
		return
	}

	accessToken, ok := tokenResponse[ACCESS_TOKEN].(string)
	if !ok {
		r.forbidden(w)
		return
	}

	newReq, err := http.NewRequest("GET", GOOGLE_USER_INFO_URL, nil)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	newReq.Header.Add("Authorization", "Bearer "+accessToken)

	userInfoResp, err := http.DefaultClient.Do(newReq)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	defer userInfoResp.Body.Close()

	var googleInfo struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&googleInfo); err != nil {
		r.serverError(w, req, err)
		return
	}

	form := entity.UserSignupForm{
		Username: googleInfo.Name,
		Email:    googleInfo.Email,
	}
	form.Password, err = pswd.GenerateRandomPassword(8)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	r.SSO(w, req, &form)
}

/*
	GITHUB LOGIN
*/

func (r *Routes) githublogin(w http.ResponseWriter, req *http.Request) {
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		r.cfg.ExternalAuth.GithubClientID, r.cfg.ExternalAuth.GithubRedirectURL,
	)

	http.Redirect(w, req, redirectURL, http.StatusTemporaryRedirect)

}

func (r *Routes) githubCallback(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	payload := url.Values{
		"client_id":     {r.cfg.ExternalAuth.GithubClientID},
		"client_secret": {r.cfg.ExternalAuth.GithubClientSecret},
		"code":          {code},
		"redirect_uri":  {r.cfg.ExternalAuth.GithubRedirectURL},
	}
	resp, err := http.PostForm(GITHUB_TOKEN_URL, payload)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	tokenResp, err := url.ParseQuery(string(body))
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	accessToken := tokenResp.Get(ACCESS_TOKEN)

	newreq, err := http.NewRequest("GET", GITHUB_USER_INFO_URL, nil)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	newreq.Header.Set("Authorization", "token "+accessToken)
	client := &http.Client{}
	resp, err = client.Do(newreq)
	if err != nil {
		r.serverError(w, req, err)
		return
	}
	defer resp.Body.Close()

	var githubInfo struct {
		Login string `json:"login"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubInfo); err != nil {
		r.serverError(w, req, err)
		return
	}

	if githubInfo.Email == "" {
		githubInfo.Email = githubInfo.Login + "@github.com"
	}
	form := entity.UserSignupForm{
		Username: githubInfo.Login,
		Email:    githubInfo.Email,
	}
	form.Password, err = pswd.GenerateRandomPassword(8)
	if err != nil {
		r.serverError(w, req, err)
		return
	}

	r.SSO(w, req, &form)
}
