package handlers

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/entity"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	RedirectURI  string
	ClientID     string
	ClientSecret string

	GithubRedirectURI  string
	GithubClientID     string
	GithubClientSecret string
)

func init() {
	file, err := os.Open(".env")
	if err != nil {
		log.Print(err)
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

	GithubRedirectURI = variables["GithubRedirect_URI"]
	GithubClientID = variables["GithubClient_id"]
	GithubClientSecret = variables["GithubClient_secret"]

}

func (r *Routes) googlelogin(w http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&redirect_uri=%s&response_type=code&scope=email profile", ClientID, RedirectURI)
	http.Redirect(w, req, url, http.StatusTemporaryRedirect)
}

func (r *Routes) googleCallback(w http.ResponseWriter, req *http.Request) {
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

	r.Sso(w, req, form)
}

func (r *Routes) Sso(w http.ResponseWriter, req *http.Request, form entity.UserSignupForm) {
	id, err := r.service.User.SaveUser(&form)

	if err != nil {

		switch {
		case errors.Is(err, entity.ErrDuplicateEmail) || errors.Is(err, entity.ErrDuplicateUsername):

			user, err := r.service.User.GetUserByEmail(form.Email)
			if err != nil {
				r.serverError(w, req, err)
				return
			}

			err = r.sesm.RenewToken(req.Context(), user.Id)
			if err != nil {
				r.serverError(w, req, err)
				return
			}

			r.sesm.PutUserID(req.Context(), user.Id)
			http.Redirect(w, req, "/", http.StatusSeeOther)
			return

		default:
			r.serverError(w, req, err)
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

// git hub

func (r *Routes) githublogin(w http.ResponseWriter, req *http.Request) {
	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s",
		GithubClientID, GithubRedirectURI,
	)

	http.Redirect(w, req, redirectURL, http.StatusTemporaryRedirect)

}

func (r *Routes) githubCallback(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")

	// Обмен code на access_token
	tokenURL := "https://github.com/login/oauth/access_token"
	payload := url.Values{
		"client_id":     {GithubClientID},
		"client_secret": {GithubClientSecret},
		"code":          {code},
		"redirect_uri":  {GithubRedirectURI},
	}
	resp, err := http.PostForm(tokenURL, payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Парсинг строки, чтобы получить токен
	tokenResp, err := url.ParseQuery(string(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken := tokenResp.Get("access_token")

	// Используйте accessToken для запросов к API GitHub от имени пользователя

	// Пример: получение информации о пользователе
	userInfoURL := "https://api.github.com/user"
	newreq, err := http.NewRequest("GET", userInfoURL, nil) // Вот здесь возникла ошибка
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newreq.Header.Set("Authorization", "token "+accessToken)
	client := &http.Client{}
	resp, err = client.Do(newreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Обработка ответа с информацией о пользователе
	// В этом примере мы просто выводим ответ на страницу

	/* userInfo, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "GitHub User Info: %s", userInfo) */

	form := entity.UserSignupForm{}

	if err := json.NewDecoder(resp.Body).Decode(&form); err != nil {
		http.Error(w, "Failed to decode user info response", http.StatusInternalServerError)
		r.serverError(w, req, err)
		return
	}

	if form.Email == "" {
		form.Email = form.Username + "@github.com"
	}

	form.Password, err = RandomPassword(8)
	if err != nil {
		http.Error(w, "Failed to generate password", http.StatusInternalServerError)
		r.serverError(w, req, err)
		return
	}

	r.Sso(w, req, form)
}
