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
		{
			name:     "Blank username",
			username: "",
			email:    validEmail,
			password: validPassword,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Blank username",
			username: validUsername,
			email:    "",
			password: validPassword,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Blank email",
			username: validUsername,
			email:    "",
			password: validPassword,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Too long username",
			username: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.",
			email:    validEmail,
			password: validPassword,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Too long email",
			username: validUsername,
			email:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.",
			password: validPassword,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Invalid email",
			username: validUsername,
			email:    "some@invalid.email.shoud.be.here@.",
			password: validPassword,
			wantCode: http.StatusBadRequest,
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
