package handlers

import (
	"fmt"
	"forum/internal/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestCommentCreate(t *testing.T) {
	r := newTestRoutes(t)
	ts := newTestServer(t, r.Register())
	defer ts.Close()

	const (
		validPostID  = "1"
		validContent = "Naaah, this one is not thaaaat great"
	)

	tests := []struct {
		name     string
		content  string
		postID   string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid comment",
			content:  validContent,
			postID:   validPostID,
			wantCode: http.StatusOK,
			wantBody: "/post/view/1",
		},
		{
			name:     "Blank content",
			content:  "",
			postID:   validPostID,
			wantCode: http.StatusBadRequest,
			wantBody: "This field cannot be blank",
		},
		{
			name:     "Invalid postID (negative)",
			content:  validContent,
			postID:   "-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Invalid postID (non-digit)",
			content:  validContent,
			postID:   "nah",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Invalid postID (non-existent)",
			content:  validContent,
			postID:   "1000000000",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("commentContent", tt.content)

			code, _, body := ts.postForm(t, fmt.Sprintf("/post/comment/%s", tt.postID), form)
			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
