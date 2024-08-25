package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"news-fetcher/internal/models"
	"news-fetcher/internal/utils"
)

type GoogleNewsAPI struct {
	APIKey string
}

func NewGoogleNewsAPI(apiKey string) *GoogleNewsAPI {
	return &GoogleNewsAPI{APIKey: apiKey}
}

func (api *GoogleNewsAPI) FetchTopHeadlinesNews(keyword string) ([]models.NewsArticleTopHeadlines, error) {
	// hardcoding 'technology' because that's our main interest
	// we're using the 'top-headlines' path instead of 'everything' because it allows us to query further by country, category, etc.
	// news are sort by 'earliest date' from the api using above path
	url := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=us&category=technology&q=%s&apiKey=%s", keyword, api.APIKey)
	resp, err := utils.MakeSecureHTTPRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("Error making http get request: %s", err)
		return nil, err
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			// Handle the error if needed, for example, log it
			log.Printf("failed to close response body: %v", err)
		}
	}()

	var result struct {
		Articles []models.NewsArticleTopHeadlines `json:"articles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding json: %s", err)
		return nil, err
	}

	log.Printf("Article Response: %s", result.Articles)

	return result.Articles, nil
}
