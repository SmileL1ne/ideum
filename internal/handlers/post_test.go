package handlers

import (
	"forum/internal/assert"
	"net/http"
	"net/url"
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
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Decimal ID",
			url:      "/post/view/1.77",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			url:      "/post/view/bruh",
			wantCode: http.StatusNotFound,
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

func TestPostCreatePost(t *testing.T) {
	r := newTestRoutes(t)
	ts := newTestServer(t, r.Register())
	defer ts.Close()

	const (
		validTitle   = "Witcher"
		validContent = "V: Wh... What are you doing??... G: Killing a monster"
	)
	var validTags = []string{"1", "2"}

	tests := []struct {
		name     string
		title    string
		content  string
		tags     []string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid post",
			title:    validTitle,
			content:  validContent,
			tags:     validTags,
			wantCode: http.StatusOK,
			wantBody: "/post/view/1",
		},
		{
			name:     "Blank title",
			title:    "",
			content:  validContent,
			tags:     validTags,
			wantCode: http.StatusBadRequest,
			wantBody: "title: This field cannot be blank",
		},
		{
			name:     "Blank content",
			title:    validTitle,
			content:  "",
			tags:     validTags,
			wantCode: http.StatusBadRequest,
			wantBody: "content: This field cannot be blank",
		},
		{
			name:     "Too long title",
			title:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.",
			content:  validContent,
			tags:     validTags,
			wantCode: http.StatusBadRequest,
			wantBody: "title: Maximum characters length exceeded",
		},
		{
			name:     "Too long content",
			title:    validTitle,
			content:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.",
			tags:     validTags,
			wantCode: http.StatusBadRequest,
			wantBody: "content: Maximum characters length exceeded",
		},
		{
			name:     "Too long content",
			title:    validTitle,
			content:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.Lorem ipsum dolor sit amet, consectetur adipiscing elit. Mauris non libero placerat, ullamcorper purus in, hendrerit tellus. Etiam quis magna sagittis lorem tincidunt gravida. Suspendisse potenti. Cras tortor nisi, suscipit id ex eu, porta blandit arcu. Aliquam lacinia lorem est, sit amet tincidunt nisi fringilla non. Nullam sit amet quam nunc. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Sed porta eget enim eu auctor. Cras ac maximus purus. Duis et tincidunt urna. Mauris nec quam sit amet massa tristique dapibus nec eu neque. Curabitur ut maximus lorem. Proin diam diam, ultricies ac condimentum nec, hendrerit at ex.",
			tags:     validTags,
			wantCode: http.StatusBadRequest,
			wantBody: "content: Maximum characters length exceeded",
		},
		{
			name:     "No tags",
			title:    validTitle,
			content:  validContent,
			wantCode: http.StatusBadRequest,
			wantBody: "tags: At least one tag should be selected",
		},
		{
			name:     "Invalid tags",
			title:    validTitle,
			content:  validContent,
			tags:     []string{"-1"},
			wantCode: http.StatusBadRequest,
			wantBody: "invalid post tags",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("title", tt.title)
			form.Add("content", tt.content)
			form["tags"] = tt.tags

			code, _, body := ts.postForm(t, "/post/create", form)
			assert.Equal(t, code, tt.wantCode)

			if tt.wantBody != "" {
				assert.StringContains(t, body, tt.wantBody)
			}
		})
	}
}
