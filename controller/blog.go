package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
	"ymmersion2/backend"
	"ymmersion2/templates"
)

func ArticlePage(w http.ResponseWriter, r *http.Request) {
	templates.Temp.ExecuteTemplate(w, "article", nil)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	templates.Temp.ExecuteTemplate(w, "index", nil)
}

func CategoriePage(w http.ResponseWriter, r *http.Request) {
	templates.Temp.ExecuteTemplate(w, "categorie", nil)
}

func ResultPage(w http.ResponseWriter, r *http.Request) {
	templates.Temp.ExecuteTemplate(w, "result", nil)
}

func RecuDatas(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)

	categorie := r.FormValue("categorie")

	titre := r.FormValue("titre")
	intro := r.FormValue("intro")
	contenu := r.FormValue("contenu")
	auteur := r.FormValue("auteur")

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	filepath := "./assets/images/" + handler.Filename
	f, _ := os.Create(filepath)
	defer f.Close()
	io.Copy(f, file)

	jsonFile, _ := ioutil.ReadFile("blog.json")
	var jsonData backend.JSONData
	json.Unmarshal(jsonFile, &jsonData)

	ids, _ := backend.GetArticleIDs("blog.json")

	rand.Seed(time.Now().UnixNano())
	var newID int

	for {
		newID = rand.Intn(8999) + 1000
		if !backend.IsIDPresent(newID, ids) {
			break
		}
	}

	nouvelArticle := backend.Article{
		Id:     newID,
		Titre:  titre,
		Image:  handler.Filename,
		Intro:  intro,
		Auteur: auteur,
		Date:   time.Now().Format("2006-01-02"),
		Body:   contenu,
	}

	backend.AddArticle(&jsonData, categorie, nouvelArticle)

	updatedData, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		fmt.Println("Erreur en convertissant les données en JSON :", err)
		return
	}

	ioutil.WriteFile("blog.json", updatedData, 0644)

	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
} // Route /new_article
