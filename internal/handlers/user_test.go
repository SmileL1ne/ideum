package handlers

import (
	"forum/internal/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestSignupPost(t *testing.T) {
	r := newTestRoutes(t)
	ts := newTestServer(t, r.Register())
	defer ts.Close()

	const (
		validUsername = "witcher"
		validEmail    = "witcher@wildHunt.com"
		validPassword = "ThisOneIsForMonsters"
	)

	tests := []struct {
		name     string
		username string
		email    string
		password string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid User",
			username: validUsername,
			email:    validEmail,
			password: validPassword,
			wantCode: http.StatusSeeOther,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("username", tt.username)
			form.Add("email", tt.email)
			form.Add("password", tt.password)

			code, _, body := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
