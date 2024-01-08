package http

import (
	"bytes"
	"html/template"
	"net/http"
)

func (r routes) newTemplateData(req *http.Request) templateData {
	return templateData{}
}

func (r routes) render(w http.ResponseWriter, req *http.Request, status int, page string, data templateData) {
	tmpl, err := template.ParseFiles(page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)
}
