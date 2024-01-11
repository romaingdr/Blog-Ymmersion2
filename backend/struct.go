package backend

var GlobalSession Session
var GlobalAccount AccountCreation

// Categorie est une structure qui stock le nom d'une catégorie et une liste d'Article
type Categorie struct {
	Name     string    `json:"name"`
	Articles []Article `json:"articles"`
}

// JSONData est une structure qui stock une liste de Categorie
type JSONData struct {
	Categories []Categorie `json:"categories"`
}

// Article est une structure qui stock toutes les données d'un article
type Article struct {
	Id     int    `json:"id"`
	Titre  string `json:"titre"`
	Image  string `json:"image"`
	Intro  string `json:"introduction"`
	Auteur string `json:"auteur"`
	Date   string `json:"date"`
	Body   string `json:"corps"`
}

// Account est une structure qui stock toutes les données d'un compte
type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	State    string `json:"state"`
	Salt     string `json:"salt"`
}

// Accounts est une structure qui stock une liste d'Account
type Accounts struct {
	Comptes []Account `json:"comptes"`
}

// Session est une structure qui stock les données d'une session
type Session struct {
	Username string
	State    string
	Mail     string
}

// IndexData est une structure qui gère les données envoyées à la page index
type IndexData struct {
	Articles   []Article
	IsLoggedIn bool
	AsAdmin    bool
}

// CategorieData est une structure qui gère les données envoyées à la page catégorie
type CategorieData struct {
	IsLoggedIn bool
	AsAdmin    bool
	Categorie  Categorie
}

// MailCodeData est une structure qui gère les données envoyées à la page de vérification mail
type MailCodeData struct {
	Success bool
	Message string
}

// ArticleData est une structure qui gère les données envoyées à la page index
type ArticleData struct {
	IsLoggedIn bool
	AsAdmin    bool
	Data       map[string]interface{}
}

// LoginStatus est une structure qui gère l'état de la connexion active sur le site
type LoginStatus struct {
	IsLoggedIn bool
	AsAdmin    bool
}

// AccountCreation est une structure qui stock des données temporaires liées à la création d'un compte
type AccountCreation struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	MailCode string
}

// AccountsCreation est une structure qui stock une liste d'AccountCreation
type AccountsCreation struct {
	Comptes []AccountCreation `json:"comptes"`
}

// RememberData est une structure qui gère l'état de la sauvegarde de session
type RememberData struct {
	Remember struct {
		Active   string `json:"Active"`
		Username string `json:"Username"`
	} `json:"remember"`
}
