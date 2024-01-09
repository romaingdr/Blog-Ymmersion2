package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"ymmersion2/backend"
	"ymmersion2/templates"
)

func AdminPage(w http.ResponseWriter, r *http.Request) {
	file, _ := ioutil.ReadFile("blog.json")

	var data backend.JSONData
	json.Unmarshal(file, &data)

	templates.Temp.ExecuteTemplate(w, "admin", data)
}

func AddArticlePage(w http.ResponseWriter, r *http.Request) {
	templates.Temp.ExecuteTemplate(w, "newarticle", nil)
}

func DeletePage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		articleID := r.FormValue("article_id")

		id, err := strconv.Atoi(articleID)
		if err != nil {
			http.Error(w, "Invalid article ID", http.StatusBadRequest)
			return
		}

		file, _ := ioutil.ReadFile("blog.json")

		var jsonData backend.JSONData
		json.Unmarshal(file, &jsonData)

		found := false
		for i := range jsonData.Categories {
			for j, article := range jsonData.Categories[i].Articles {
				if article.Id == id {
					jsonData.Categories[i].Articles = append(jsonData.Categories[i].Articles[:j], jsonData.Categories[i].Articles[j+1:]...)
					found = true
					break
				}
			}
			if found {
				break
			}
		}

		newData, _ := json.MarshalIndent(jsonData, "", "  ")
		ioutil.WriteFile("blog.json", newData, 0644)

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	err := templates.Temp.ExecuteTemplate(w, "erreur", nil)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
}
