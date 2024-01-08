package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {

	css := http.FileServer(http.Dir("./client/style"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	tmpl, _ := template.ParseGlob("./templates/*.gohtml")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index", nil)
		if err != nil {
			log.Println(err)
			http.Error(w, "Erreur : ", http.StatusInternalServerError)
		}
	})

	// Démarrage du serveur
	log.Println("[✅] Serveur lancé !")
	fmt.Println("[🌐] http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
