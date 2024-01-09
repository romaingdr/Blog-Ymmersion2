package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
	"ymmersion2/backend"
	"ymmersion2/templates"
)

func ArticlePage(w http.ResponseWriter, r *http.Request) {
	queryID := r.URL.Query().Get("id")
	articleID, err := strconv.Atoi(queryID)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var jsonData backend.JSONData

	jsonDataFile, err := ioutil.ReadFile("blog.json")
	err = json.Unmarshal(jsonDataFile, &jsonData)

	var foundArticle *backend.Article
	for _, category := range jsonData.Categories {
		for _, article := range category.Articles {
			if article.Id == articleID {
				foundArticle = &article
				break
			}
		}
		if foundArticle != nil {
			break
		}
	}

	if foundArticle == nil {
		http.Error(w, "Article not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Article": foundArticle,
	}

	templates.Temp.ExecuteTemplate(w, "article", data)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	content, err := os.ReadFile("blog.json")
	if err != nil {
		fmt.Println("Erreur dans la lecture du json : ", err)
	}

	var result backend.JSONData

	err = json.Unmarshal(content, &result)
	if err != nil {
		fmt.Println("Erreur > ", err.Error())
	}

	randomArticles := getRandomArticles(result, 10)

	templates.Temp.ExecuteTemplate(w, "index", randomArticles)
}

// Prends un Article au hasard
func getRandomArticles(data backend.JSONData, count int) []backend.Article {
	var randomArticles []backend.Article
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count && i < len(data.Categories); i++ {
		category := data.Categories[i]
		if len(category.Articles) > 0 {
			randomIndex := rand.Intn(len(category.Articles))
			randomArticles = append(randomArticles, category.Articles[randomIndex])
		}
	}

	return randomArticles
}

func CategoriePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")

	content, err := os.ReadFile("blog.json")
	if err != nil {
		fmt.Println("Erreur dans la lecture du json : ", err)
	}

	var result backend.JSONData

	err = json.Unmarshal(content, &result)
	if err != nil {
		fmt.Println("Erreur > ", err.Error())
	}

	var Data backend.Categorie

	urlStr := r.URL.RawQuery
	fmt.Println(urlStr)
	switch urlStr {
	case "categorie=esport":
		Data = result.Categories[0]
	case "categorie=nouveautes":
		Data = result.Categories[1]
	case "categorie=presentations":
		Data = result.Categories[2]
	default:
		tmpl.ExecuteTemplate(w, "erreur", nil)
	}

	fmt.Println(Data)

	tmpl.ExecuteTemplate(w, "categorie", Data)
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
