package http

import (
	"forum/internal/entity"
	"forum/web"
	"html/template"
	"io/fs"
	"path/filepath"
)

type templateData struct {
	Posts []entity.Post
	Post  entity.Post
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
			page,
		}

		ts, err := template.ParseFS(web.Files, files...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
