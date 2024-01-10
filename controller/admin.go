package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"ymmersion2/backend"
	"ymmersion2/templates"
)

var GlobalSession backend.Session

func AdminPage(w http.ResponseWriter, r *http.Request) {
	if backend.GetSession().State != "admin" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}
	file, _ := ioutil.ReadFile("blog.json")

	var data backend.JSONData
	json.Unmarshal(file, &data)

	templates.Temp.ExecuteTemplate(w, "admin", data)
}

func AddArticlePage(w http.ResponseWriter, r *http.Request) {
	if backend.GetSession().State != "admin" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}
	templates.Temp.ExecuteTemplate(w, "newarticle", nil)
}

func DeletePage(w http.ResponseWriter, r *http.Request) {
	if backend.GetSession().State != "admin" {
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
		return
	}
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

func LoginPage(w http.ResponseWriter, r *http.Request) {
	errMessage := r.URL.Query().Get("error")
	templates.Temp.ExecuteTemplate(w, "login", errMessage)
}

func GetCreds(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Erreur lors de la récupération des données du formulaire :", err)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	file, err := ioutil.ReadFile("accounts.json")
	if err != nil {
		fmt.Println("Erreur de lecture du fichier JSON :", err)
		http.Redirect(w, r, "/login?error=Erreur de lecture du fichier JSON", http.StatusSeeOther)
		return
	}

	var accounts backend.Accounts
	err = json.Unmarshal(file, &accounts)
	if err != nil {
		fmt.Println("Erreur lors de la conversion du JSON :", err)
		http.Redirect(w, r, "/login?error=Erreur lors de la conversion du JSON", http.StatusSeeOther)
		return
	}

	valid := false
	for _, account := range accounts.Comptes {
		if account.Username == username && account.Password == password {
			valid = true
			break
		}
	}

	if valid {
		session := backend.Session{Username: username, State: backend.GetAccountState(username)}
		backend.SetSession(session)
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login?error=Identifiants invalides", http.StatusSeeOther)
	}
}

func Deconnexion(w http.ResponseWriter, r *http.Request) {
	backend.ClearSession()
	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}
