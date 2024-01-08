package backend

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
