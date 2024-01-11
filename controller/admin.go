package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
	"ymmersion2/backend"
	"ymmersion2/templates"
)

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

	mail := r.Form.Get("email")
	password := r.Form.Get("password")
	remember := r.FormValue("remember")

	file, _ := ioutil.ReadFile("accounts.json")

	var accounts backend.Accounts
	json.Unmarshal(file, &accounts)

	valid := false
	for _, account := range accounts.Comptes {
		if account.Email == mail || account.Username == mail {
			if backend.HashPassword(password, account.Salt) == account.Password {
				fmt.Println("here")
				valid = true
				break
			}
		}
	}

	if valid {
		username := backend.GetUsernameByEmail(mail)
		session := backend.Session{Username: username, State: backend.GetAccountState(username), Mail: mail}
		if remember == "on" {
			backend.SetRememberActive(username, "rememberSession.json")
		}
		fmt.Println(session)
		backend.SetSession(session)
		http.Redirect(w, r, "/accueil", http.StatusSeeOther)
	} else {
		fmt.Println("Identifiants incorrects")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func Deconnexion(w http.ResponseWriter, r *http.Request) {
	backend.ClearSession()
	backend.ClearRemember("rememberSession.json")
	http.Redirect(w, r, "/accueil", http.StatusSeeOther)
}

func MailVerifPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		data := backend.MailCodeData{
			Success: false,
			Message: "Le code est incorrect !",
		}
		templates.Temp.ExecuteTemplate(w, "mailverif", data)
		return
	}

	r.ParseForm()

	emailDestinataire := r.FormValue("email")
	username := r.FormValue("username")
	passwordAccount := r.FormValue("password")

	for _, element := range backend.GetEmailsFromJSON("accounts.json") {
		if element == emailDestinataire {
			fmt.Println("Mail déjà utilisé")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	if len(username) < 5 {
		fmt.Println("Nom d'utilisateur trop court")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	} else {
		for _, element := range backend.GetUsersFromJSON("accounts.json") {
			if element == username {
				fmt.Println("Nom d'utilisateur déjà utilisé")
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
		}
	}

	if len(passwordAccount) < 8 {
		fmt.Println("Mot de passe trop court")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	rand.Seed(time.Now().UnixNano())

	codeMail := rand.Intn(89999) + 10000
	codeMailString := strconv.Itoa(codeMail)

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	email := "octogamesverify@gmail.com"
	password := "womp qoly znmc krqe"

	to := []string{emailDestinataire}
	subject := "Code de vérification Octo Games"
	body := "Bonjour " + username + ",\nVoici votre code de vérification : " + codeMailString

	auth := smtp.PlainAuth("", email, password, smtpHost)

	msg := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	err := smtp.SendMail(smtpHost+":"+strconv.Itoa(smtpPort), auth, email, to, msg)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("E-mail envoyé avec succès!")

	tempaccount := backend.AccountCreation{Username: username, Email: emailDestinataire, Password: passwordAccount, MailCode: codeMailString}
	backend.GlobalAccount = tempaccount

	data := backend.MailCodeData{
		Success: true,
		Message: "",
	}

	templates.Temp.ExecuteTemplate(w, "mailverif", data)
}

func VerifCode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	codeRecu := r.FormValue("verificationCode")
	codeEnvoye := backend.GlobalAccount.MailCode

	fmt.Println(codeRecu)
	fmt.Println(codeEnvoye)
	fmt.Println("Verification du code")

	if codeRecu == codeEnvoye {
		http.Redirect(w, r, "/success_code", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/mail_verif", http.StatusSeeOther)
}

func SuccessPage(w http.ResponseWriter, r *http.Request) {
	backend.AddAccountToFile(backend.GlobalAccount, "accounts.json")
	backend.ClearAccount()
	fmt.Println(backend.GlobalAccount)
	templates.Temp.ExecuteTemplate(w, "success", nil)
}
