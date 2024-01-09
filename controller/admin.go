package controller

import (
	"net/http"
	"ymmersion2/templates"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	templates.Temp.ExecuteTemplate(w, "admin", nil)
}

func AddArticlePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/new_article" {
		NotFoundPage(w, r, http.StatusNotFound)
		return
	}

	templates.Temp.ExecuteTemplate(w, "newarticle", nil)
}
