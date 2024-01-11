package backend

var GlobalSession Session
var GlobalAccount AccountCreation

type Categorie struct {
	Name     string    `json:"name"`
	Articles []Article `json:"articles"`
}

type JSONData struct {
	Categories []Categorie `json:"categories"`
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

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	State    string `json:"state"`
	Salt     string `json:"salt"`
}

type Accounts struct {
	Comptes []Account `json:"comptes"`
}

type AccountsCreation struct {
	Comptes []AccountCreation `json:"comptes"`
}

type Session struct {
	Username string
	State    string
	Mail     string
}

type IndexData struct {
	Articles   []Article
	IsLoggedIn bool
	AsAdmin    bool
}

type CategorieData struct {
	IsLoggedIn bool
	AsAdmin    bool
	Categorie  Categorie
}

type ArticleData struct {
	IsLoggedIn bool
	AsAdmin    bool
	Data       map[string]interface{}
}

type LoginStatus struct {
	IsLoggedIn bool
	AsAdmin    bool
}

type AccountCreation struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	MailCode string
}

type MailCodeData struct {
	Success bool
	Message string
}
