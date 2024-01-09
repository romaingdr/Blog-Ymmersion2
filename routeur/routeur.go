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
	http.HandleFunc("/delete", controller.DeletePage)
	http.HandleFunc("/mentions_legales", controller.Mentions)

	http.HandleFunc("/", controller.DefaultHandler)

	// Démarrage du serveur
	log.Println("[✅] Serveur lancé !")
	fmt.Println("[🌐] http://localhost:8080/admin")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
