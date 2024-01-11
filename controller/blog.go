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
	"time"
	"ymmersion2/backend"
	"ymmersion2/templates"
)

func ArticlePage(w http.ResponseWriter, r *http.Request) {
	session := backend.GetSession() != backend.Session{}
	isAdmin := backend.IsAdmin()

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
		templates.Temp.ExecuteTemplate(w, "erreur", nil)
	}

	data := map[string]interface{}{
		"Article": foundArticle,
	}

	articleData := backend.ArticleData{
		IsLoggedIn: session,
		AsAdmin:    isAdmin,
		Data:       data,
	}

	templates.Temp.ExecuteTemplate(w, "article", articleData)
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	session := backend.GetSession() != backend.Session{}
	isAdmin := backend.IsAdmin()
	content, err := os.ReadFile("blog.json")
	if err != nil {
		fmt.Println("Erreur dans la lecture du json : ", err)
	}

	var result backend.JSONData

	err = json.Unmarshal(content, &result)
	if err != nil {
		fmt.Println("Erreur > ", err.Error())
	}

	randomArticles := backend.GetRandomArticles(result)

	data := backend.IndexData{
		Articles:   randomArticles,
		IsLoggedIn: session,
		AsAdmin:    isAdmin,
	}

	fmt.Println(session)
	fmt.Println(isAdmin)
	templates.Temp.ExecuteTemplate(w, "index", data)
}

func CategoriePage(w http.ResponseWriter, r *http.Request) {
	session := backend.GetSession() != backend.Session{}
	isAdmin := backend.IsAdmin()

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
	switch urlStr {
	case "categorie=esport":
		Data = result.Categories[0]
	case "categorie=nouveautes":
		Data = result.Categories[1]
	case "categorie=presentations":
		Data = result.Categories[2]
	default:
		http.Redirect(w, r, "/error", http.StatusSeeOther)
	}

	categorieData := backend.CategorieData{
		IsLoggedIn: session,
		AsAdmin:    isAdmin,
		Categorie:  Data,
	}

	templates.Temp.ExecuteTemplate(w, "categorie", categorieData)
}

func ResultPage(w http.ResponseWriter, r *http.Request) {
	session := backend.GetSession() != backend.Session{}
	isAdmin := backend.IsAdmin()
	recherche := r.URL.Query().Get("content")
	var jsonData backend.JSONData

	file, err := ioutil.ReadFile("blog.json")
	if err != nil {
		http.Error(w, "Impossible de charger les données", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du fichier JSON", http.StatusInternalServerError)
		return
	}

	var resultArticles []backend.Article

	for _, cat := range jsonData.Categories {
		for _, article := range cat.Articles {
			if backend.TitleContains(article.Titre, recherche) {
				resultArticles = append(resultArticles, article)
			}
		}
	}

	data := backend.IndexData{
		Articles:   resultArticles,
		IsLoggedIn: session,
		AsAdmin:    isAdmin,
	}

	templates.Temp.ExecuteTemplate(w, "result", data)
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

	http.Redirect(w, r, "/new_article", http.StatusSeeOther)
}

func Mentions(w http.ResponseWriter, r *http.Request) {
	data := backend.LoginStatus{IsLoggedIn: backend.GetSession() != backend.Session{}, AsAdmin: backend.IsAdmin()}

	templates.Temp.ExecuteTemplate(w, "mentions", data)
}

func Repartition(w http.ResponseWriter, r *http.Request) {
	data := backend.LoginStatus{IsLoggedIn: backend.GetSession() != backend.Session{}, AsAdmin: backend.IsAdmin()}

	templates.Temp.ExecuteTemplate(w, "repartition", data)
}

func Explication(w http.ResponseWriter, r *http.Request) {
	data := backend.LoginStatus{IsLoggedIn: backend.GetSession() != backend.Session{}, AsAdmin: backend.IsAdmin()}

	templates.Temp.ExecuteTemplate(w, "explication", data)
}
