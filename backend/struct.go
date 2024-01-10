package backend

var GlobalSession Session

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
	State    string `json:"state"`
}

type Accounts struct {
	Comptes []Account `json:"comptes"`
}

type Session struct {
	Username string
	State    string
}

type IndexData struct {
	Articles   []Article
	IsLoggedIn bool
}

type CategorieData struct {
	IsLoggedIn bool
	Categorie  Categorie
}

type ArticleData struct {
	IsLoggedIn bool
	Data       map[string]interface{}
}

type LoginStatus struct {
	IsLoggedIn bool
}
