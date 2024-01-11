package backend

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

// IsIDPresent vérifie si un id est présent dans une liste d'id
func IsIDPresent(id int, ids []int) bool {
	for _, existingID := range ids {
		if existingID == id {
			return true
		}
	}
	return false
}

// GetArticleIDs récupère tous les id des articles du fichier blog.json
func GetArticleIDs(filename string) ([]int, error) {
	var data map[string]interface{}

	raw, _ := ioutil.ReadFile(filename)

	json.Unmarshal(raw, &data)

	var articleIDs []int
	categories, ok := data["categories"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("Champ 'categories' non trouvé ou incorrect")
	}

	for _, category := range categories {
		cat, ok := category.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Structure de catégorie incorrecte")
		}

		articles, ok := cat["articles"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("Champ 'articles' non trouvé ou incorrect")
		}

		for _, article := range articles {
			art, ok := article.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("Structure d'article incorrecte")
			}

			id, ok := art["id"].(float64)
			if !ok {
				return nil, fmt.Errorf("Champ 'id' non trouvé ou incorrect")
			}

			articleIDs = append(articleIDs, int(id))
		}
	}

	return articleIDs, nil
}

// TitleContains vérifie si une chaine de caractère substr est contenue dans une chaine s
func TitleContains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// AddArticle prend en argument une structure json, un nom de catégorie et un article puis ajoute l'article dans la bonne catégorie
func AddArticle(jsonData *JSONData, categoryName string, article Article) error {
	for i := range jsonData.Categories {
		if jsonData.Categories[i].Name == categoryName {
			jsonData.Categories[i].Articles = append(jsonData.Categories[i].Articles, article)
			return nil
		}
	}
	return fmt.Errorf("category '%s' not found", categoryName)
}

// GetAllArticles récupère tous les articles contenus dans une structure JSONData
func GetAllArticles(jsonData JSONData) []Article {
	var allArticles []Article
	for _, categorie := range jsonData.Categories {
		allArticles = append(allArticles, categorie.Articles...)
	}

	return allArticles
}

// GetRandomArticles récupère 10 articles aléatoires parmi une liste complète d'articles
func GetRandomArticles(jsonData JSONData) []Article {
	rand.Seed(time.Now().UnixNano())

	allArticles := GetAllArticles(jsonData)

	if len(allArticles) <= 10 {
		return allArticles
	}

	rand.Shuffle(len(allArticles), func(i, j int) {
		allArticles[i], allArticles[j] = allArticles[j], allArticles[i]
	})
	return allArticles[:10]
}

// GetAccountState récupère le statut d'un utilisateur par son pseudonyme dans le fichier accounts.json
func GetAccountState(username string) string {
	file, _ := ioutil.ReadFile("accounts.json")

	var accounts Accounts
	json.Unmarshal(file, &accounts)

	for _, account := range accounts.Comptes {
		if account.Username == username {
			return account.State
		}
	}

	return ""
}

// GetAccountMail récupère le mail d'un utilisateur par son pseudonyme dans le fichier accounts.json
func GetAccountMail(username string) string {
	file, _ := ioutil.ReadFile("accounts.json")

	var accounts Accounts
	json.Unmarshal(file, &accounts)

	for _, account := range accounts.Comptes {
		if account.Username == username {
			return account.Email
		}
	}

	return ""
}

// SetSession paramètre la session utilisateur globale active sur le site
func SetSession(session Session) {
	GlobalSession = session
}

// GetSession renvoie la session utilisateur globale active sur le site
func GetSession() Session {
	return GlobalSession
}

// IsAdmin renvoie si la session utilisateur en cours est une session administrateur
func IsAdmin() bool {
	return GlobalSession.State == "admin"
}

// ClearSession vide la session en cours sur le site
func ClearSession() {
	GlobalSession = Session{}
}

// ClearAccount vide la variable temporaire GlobalAccount pour la création de compte
func ClearAccount() {
	GlobalAccount = AccountCreation{}
}

// AddAccountToFile ajoute une variable temporaire de création de compte à un fichier json
func AddAccountToFile(account AccountCreation, filePath string) error {
	jsonFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Erreur lors de la lecture du fichier JSON : %v", err)
	}

	var data map[string][]map[string]interface{}
	err = json.Unmarshal(jsonFile, &data)
	if err != nil {
		return fmt.Errorf("Erreur lors du parsing du JSON : %v", err)
	}

	salt, err := GenerateSalt()
	hashedPassword := HashPassword(account.Password, salt)

	newAccount := map[string]interface{}{
		"username": account.Username,
		"email":    account.Email,
		"password": hashedPassword,
		"state":    "member",
		"salt":     salt,
	}

	accounts, ok := data["comptes"]
	if !ok {
		accounts = make([]map[string]interface{}, 0)
	}
	data["comptes"] = append(accounts, newAccount)

	newJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Erreur lors de la conversion en JSON : %v", err)
	}

	err = ioutil.WriteFile(filePath, newJSON, 0644)
	if err != nil {
		return fmt.Errorf("Erreur lors de l'écriture dans le fichier JSON : %v", err)
	}

	return nil
}

// GetUsernameByEmail permet de récupérer le nom d'utilisateur d'un compte par son mail dans la base de données
func GetUsernameByEmail(emailToFind string) string {
	jsonFile, _ := ioutil.ReadFile("accounts.json")

	var data Accounts
	json.Unmarshal(jsonFile, &data)

	for _, account := range data.Comptes {
		if account.Email == emailToFind {
			return account.Username
		}
	}

	return ""
}

// GetEmailsFromJSON permet de récupérer la liste des emails présents dans le fichier json
func GetEmailsFromJSON(filePath string) []string {
	fileContent, _ := ioutil.ReadFile(filePath)

	var comptes Accounts

	json.Unmarshal(fileContent, &comptes)

	var emails []string
	for _, compte := range comptes.Comptes {
		emails = append(emails, compte.Email)
	}

	return emails
}

// GetUsersFromJSON permet de récupérer la liste des noms d'utilisateurs présents dans le fichier json
func GetUsersFromJSON(filePath string) []string {
	fileContent, _ := ioutil.ReadFile(filePath)

	var comptes Accounts

	json.Unmarshal(fileContent, &comptes)

	var usernames []string
	for _, compte := range comptes.Comptes {
		usernames = append(usernames, compte.Username)
	}

	return usernames
}

// GenerateSalt permet de générer un sel afin de hasher un mot de passe
func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(salt), nil
}

// HashPassword permet de hasher un mot de passe avec un sel prédéfini
func HashPassword(password string, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

// CheckRememberStatus permet de lire un fichier json afin de savoir si une session a été sauvegardé par l'utilisateur
func CheckRememberStatus(filename string) (bool, string) {
	content, _ := ioutil.ReadFile(filename)

	var data RememberData
	json.Unmarshal(content, &data)

	if data.Remember.Active == "True" {
		return true, data.Remember.Username
	}

	return false, ""
}

// SetRememberActive permet d'ajouter une sauvegarde de session avec le nom d'utilisateur de la session active
func SetRememberActive(username string, filename string) error {
	content, _ := ioutil.ReadFile(filename)

	var data RememberData
	json.Unmarshal(content, &data)

	data.Remember.Active = "True"
	data.Remember.Username = username

	newContent, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Erreur lors de la création du nouveau contenu JSON : %v", err)
	}

	ioutil.WriteFile(filename, newContent, 0644)
	return nil
}

// ClearRemember permet de supprimer la sauvegarde de session active
func ClearRemember(filename string) error {
	content, _ := ioutil.ReadFile(filename)

	var data RememberData
	json.Unmarshal(content, &data)

	data.Remember.Active = "False"
	data.Remember.Username = ""

	newContent, _ := json.MarshalIndent(data, "", "  ")

	ioutil.WriteFile(filename, newContent, 0644)

	return nil
}
