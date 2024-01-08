package routeur

import (
	"fmt"
	"log"
	"net/http"
	"ymmersion2/controller"
)

func Initserv() {

	css := http.FileServer(http.Dir("./assets/"))
	http.Handle("/static/", http.StripPrefix("/static/", css))

	http.HandleFunc("/accueil", controller.IndexPage)
	http.HandleFunc("/article", controller.ArticlePage)
	http.HandleFunc("/categorie", controller.CategoriePage)
	http.HandleFunc("/admin", controller.AdminPage)
	http.HandleFunc("/result", controller.ResultPage)
	http.HandleFunc("/new_article", controller.AddArticlePage)
	http.HandleFunc("/submit", controller.RecuDatas)

	// Démarrage du serveur
	log.Println("[✅] Serveur lancé !")
	fmt.Println("[🌐] http://localhost:8080/new_article")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
