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

func IsIDPresent(id int, ids []int) bool {
	for _, existingID := range ids {
		if existingID == id {
			return true
		}
	}
	return false
}

func GetArticleIDs(filename string) ([]int, error) {
	var data map[string]interface{}

	// Lecture du fichier JSON
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Conversion du JSON en une carte générique
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return nil, err
	}

	// Récupération des IDs des articles
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

func TitleContains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func AddArticle(jsonData *JSONData, categoryName string, article Article) error {
	for i := range jsonData.Categories {
		if jsonData.Categories[i].Name == categoryName {
			jsonData.Categories[i].Articles = append(jsonData.Categories[i].Articles, article)
			return nil
		}
	}
	return fmt.Errorf("category '%s' not found", categoryName)
}

func GetAllArticles(jsonData JSONData) []Article {
	var allArticles []Article
	for _, categorie := range jsonData.Categories {
		allArticles = append(allArticles, categorie.Articles...)
	}

	return allArticles
}

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

func SetSession(session Session) {
	GlobalSession = session
}

func GetSession() Session {
	return GlobalSession
}

func ClearSession() {
	GlobalSession = Session{}
}

func ClearAccount() {
	GlobalAccount = AccountCreation{}
}

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
		"state":    "admin",
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

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(salt), nil
}

func HashPassword(password string, salt string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password + salt))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
