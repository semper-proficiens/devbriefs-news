package api

import (
	"encoding/json"
	"fmt"
	"internal/models"
	"internal/utils"
	"net/http"
)

type GoogleNewsAPI struct {
	APIKey string
}

func NewGoogleNewsAPI(apiKey string) *GoogleNewsAPI {
	return &GoogleNewsAPI{APIKey: apiKey}
}

func (api *GoogleNewsAPI) FetchNews(category, keyword string) ([]models.NewsArticle, error) {
	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&category=%s&apiKey=%s", keyword, category, api.APIKey)
	resp, err := utils.MakeHTTPRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Articles []models.NewsArticle `json:"articles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Articles, nil
}
