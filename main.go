package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type JsonData struct {
	Categories []Categorie `json:"categories"`
}

type Categorie struct {
	Name     string    `json:"name"`
	Articles []Article `json:"articles"`
}

type Article struct {
	Id     int    `json:"id"`
	Titre  string `json:"titre"`
	Image  string `json:"image"`
	Intro  string `json:"introduction"`
	Auteur string `json:"auteur"`
	Date   string `json:"date"`
	Body   string `json:"corps"`
}

func articlePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")
	tmpl.ExecuteTemplate(w, "article", nil)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")
	tmpl.ExecuteTemplate(w, "index", nil)
}

func categoriePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")

	content, err := os.ReadFile("blog.json")
	if err != nil {
		fmt.Println("Erreur dans la lecture du json : ", err)
	}

	var result JsonData

	err = json.Unmarshal(content, &result)
	if err != nil {
		fmt.Println("Erreur > ", err.Error())
	}

	var Data Categorie

	switch urlStr := r.URL.RawQuery[9:]; urlStr {
	case "esport":
		Data = result.Categories[0]
	case "nouveautes":
		Data = result.Categories[1]
	case "presentations":
		Data = result.Categories[2]
	default:
		tmpl.ExecuteTemplate(w, "erreur", nil)
	}

	fmt.Println(Data)

	tmpl.ExecuteTemplate(w, "categorie", Data)
}

func adminPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")
	tmpl.ExecuteTemplate(w, "admin", nil)
}

func resultPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")
	tmpl.ExecuteTemplate(w, "result", nil)
}

func addArticlePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseGlob("./templates/*.gohtml")
	tmpl.ExecuteTemplate(w, "newarticle", nil)
}

func main() {

	css := http.FileServer(http.Dir("./client/style"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	http.HandleFunc("/accueil", indexPage)
	http.HandleFunc("/article", articlePage)
	http.HandleFunc("/categorie", categoriePage)
	http.HandleFunc("/admin", adminPage)
	http.HandleFunc("/result", resultPage)
	http.HandleFunc("/new_article", addArticlePage)

	// Démarrage du serveur
	log.Println("[✅] Serveur lancé !")
	fmt.Println("[🌐] http://localhost:8080/accueil")
	fmt.Println("[🌐] http://localhost:8080/article")
	fmt.Println("[🌐] http://localhost:8080/categorie")
	fmt.Println("[🌐] http://localhost:8080/admin")
	fmt.Println("[🌐] http://localhost:8080/result")
	fmt.Println("[🌐] http://localhost:8080/new_article")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
