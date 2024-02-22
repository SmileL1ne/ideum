package main

import (
	"bytes"
	"forum/internal/entity/mocks"
	"forum/pkg/sesm"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
)

const (
	testUsername        = "nah"
	testPassword        = "nahnahnah"
	sessionNameInCookie = "session"
)

func newTestRoutes(t *testing.T) *routes {
	repos := mocks.NewReposMock()
	services := mocks.NewServicesMock(repos)

	tempCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	sesm := sesm.New()

	return &routes{
		service:   services,
		tempCache: tempCache,
		sesm:      sesm,
		logger:    log.New(io.Discard, "", 0),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, url string) (int, http.Header, string) {
	rs, err := ts.Client().Get(ts.URL + url)
	if err != nil {
		t.Fatal(err)
	}
	defer rs.Body.Close()

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) postForm(t *testing.T, url string, form url.Values) (int, http.Header, string) {
	req, err := http.NewRequest("POST", ts.URL+url, bytes.NewBufferString(form.Encode()))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cookie := &http.Cookie{
		Name: sessionNameInCookie,
	}

	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)

	return rs.StatusCode, rs.Header, string(body)
}

func (ts *testServer) logIn(t *testing.T) string {
	form := url.Values{}
	form.Add("identifier", testUsername)
	form.Add("password", testPassword)

	rs, err := ts.Client().PostForm(ts.URL+"/user/login", form)
	if err != nil {
		t.Fatal(err)
	}

}
