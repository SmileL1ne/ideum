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
			name:     "Valid signup",
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
			wantBody: "username: This field cannot be blank",
		},
		{
			name:     "Too long username",
			username: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.",
			email:    validEmail,
			password: validPassword,
			wantCode: http.StatusBadRequest,
			wantBody: "username: Maximum characters length exceeded - ",
		},
		{
			name:     "Invalid username (non-ascii)",
			username: "нееееееееееее",
			email:    validEmail,
			password: validPassword,
			wantCode: http.StatusBadRequest,
			wantBody: "username: Only valid characters (ascii standard) should be included",
		},
		{
			name:     "Duplicate username",
			username: "satoru",
			email:    validEmail,
			password: validPassword,
			wantCode: http.StatusBadRequest,
			wantBody: "username: Username is already in use",
		},
		{
			name:     "Blank email",
			username: validUsername,
			email:    "",
			password: validPassword,
			wantCode: http.StatusBadRequest,
			wantBody: "email: This field cannot be blank",
		},
		{
			name:     "Too long email",
			username: validUsername,
			email:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.",
			password: validPassword,
			wantCode: http.StatusBadRequest,
			wantBody: "email: Maximum characters length exceeded - ",
		},
		{
			name:     "Invalid email",
			username: validUsername,
			email:    "some@invalid.email.shoud.be.here@.",
			password: validPassword,
			wantCode: http.StatusBadRequest,
			wantBody: "email: Invalid email address",
		},
		{
			name:     "Invalid email (non-ascii)",
			username: validUsername,
			email:    "недействительный@мейл.ру",
			password: validPassword,
			wantCode: http.StatusBadRequest,
			wantBody: "email: Invalid email address",
		},
		{
			name:     "Duplicate email",
			username: validUsername,
			email:    "satoru@gmail.com",
			password: validPassword,
			wantCode: http.StatusBadRequest,
			wantBody: "email: Email address is already in use",
		},
		{
			name:     "Blank password",
			username: validUsername,
			email:    validEmail,
			password: "",
			wantCode: http.StatusBadRequest,
			wantBody: "password: This field cannot be blank",
		},
		{
			name:     "Short password",
			username: validUsername,
			email:    validEmail,
			password: "a",
			wantCode: http.StatusBadRequest,
			wantBody: "password: Minimum length for password: ",
		},
		{
			name:     "Too long password",
			username: validUsername,
			email:    validEmail,
			password: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam porttitor laoreet nisi eget molestie. Morbi vestibulum enim nec pharetra mattis. Etiam iaculis consequat risus, et facilisis elit venenatis ac. Suspendisse at consectetur nibh, quis interdum leo. Ut convallis eget justo vitae condimentum. Vivamus justo mauris, iaculis vitae ex nec, vehicula blandit est. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Praesent aliquet fermentum turpis nec rutrum.",
			wantCode: http.StatusBadRequest,
			wantBody: "password: Maximum characters length exceeded - ",
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

func TestUserLoginPost(t *testing.T) {
	r := newTestRoutes(t)
	ts := newTestServer(t, r.Register())
	defer ts.Close()

	const (
		validUsername = "yuta"
		validEmail    = "yuta@gmail.com"
		validPassword = "all-encompassing unequivocal love"
	)

	tests := []struct {
		name       string
		identifier string
		password   string
		wantCode   int
		wantBody   string
	}{
		{
			name:       "Valid login (username)",
			identifier: validUsername,
			password:   validPassword,
			wantCode:   http.StatusSeeOther,
		},
		{
			name:       "Valid login (email)",
			identifier: validEmail,
			password:   validPassword,
			wantCode:   http.StatusSeeOther,
		},
		{
			name:       "Blank identifier",
			identifier: "",
			password:   validPassword,
			wantCode:   http.StatusBadRequest,
			wantBody:   "identifier: This field cannot be blank",
		},
		{
			name:       "Invalid identifier (incorrect email)",
			identifier: "naaaaaaaah@gmail.com@gmail.com",
			password:   validPassword,
			wantCode:   http.StatusBadRequest,
			wantBody:   "Email or password is incorrect",
		},
		{
			name:       "Invalid identifier (non-existent username)",
			identifier: "I_dont't_exist",
			password:   validPassword,
			wantCode:   http.StatusBadRequest,
			wantBody:   "Email or password is incorrect",
		},
		{
			name:       "Invalid identifier (non-existent email)",
			identifier: "dontexist@gmail.com",
			password:   validPassword,
			wantCode:   http.StatusBadRequest,
			wantBody:   "Email or password is incorrect",
		},
		{
			name:       "Blank password",
			identifier: validEmail,
			password:   "",
			wantCode:   http.StatusBadRequest,
			wantBody:   "password: This field cannot be blank",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("identifier", tt.identifier)
			form.Add("password", tt.password)

			code, _, body := ts.postForm(t, "/user/login", form)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
