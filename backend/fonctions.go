package backend

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func IsIDPresent(id int, ids []int) bool {
	for _, existingID := range ids {
		if existingID == id {
			return true
		}
	}
	return false
} // Route /new_article

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
} // Route /new_article

func AddArticle(jsonData *JSONData, categoryName string, article Article) error {
	for i := range jsonData.Categories {
		if jsonData.Categories[i].Name == categoryName {
			jsonData.Categories[i].Articles = append(jsonData.Categories[i].Articles, article)
			return nil
		}
	}
	return fmt.Errorf("category '%s' not found", categoryName)
} // Route /new_article
