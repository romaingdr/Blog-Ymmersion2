package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

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
	tmpl.ExecuteTemplate(w, "categorie", nil)
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
