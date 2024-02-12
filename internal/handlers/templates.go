package handlers

import (
	"forum/internal/entity"
	"forum/web"
	"html/template"
	"io/fs"
	"path/filepath"
	"strings"
)

type Models struct {
	Posts    []entity.PostView
	Post     entity.PostView
	Comments []entity.CommentView
	Tags     []entity.TagEntity
}

type templateData struct {
	Models          Models
	Flash           string
	Form            interface{}
	IsAuthenticated bool
}

/* testing temp functions */

var fm = template.FuncMap{
	"low": strings.ToLower, // for icons
	"rev": reverse,
}

// sort posts by date
func reverse(slice []entity.PostView) []entity.PostView {
	length := len(slice)
	reversed := make([]entity.PostView, length)
	for i, v := range slice {
		reversed[length-i-1] = v
	}
	return reversed
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(web.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"html/base.html",
			"html/partials/nav.html",
			"html/partials/topics.html",
			"html/partials/userbar.html",
			page,
		}

		ts, err := template.New("").Funcs(fm).ParseFS(web.Files, files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
