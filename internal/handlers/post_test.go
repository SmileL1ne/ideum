package handlers

import (
	"forum/internal/assert"
	"net/http"
	"testing"
)

func TestPostView(t *testing.T) {
	r := newTestRoutes(t)

	ts := newTestServer(t, r.Register())
	defer ts.Close()

	tests := []struct {
		name     string
		url      string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			url:      "/post/view/1",
			wantCode: http.StatusOK,
			wantBody: "You can do it!",
		},
		{
			name:     "Non-existent ID",
			url:      "/post/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			url:      "/post/view/-1",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Decimal ID",
			url:      "/post/view/1.77",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "String ID",
			url:      "/post/view/bruh",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty ID",
			url:      "/snippet/view",
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, body := ts.get(t, tt.url)

			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
