package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"news-fetcher/internal/models"
	"news-fetcher/internal/utils"
)

const (
	hackingQuery = "`data breach OR hack`"
)

type GoogleNewsAPI struct {
	APIKey string
}

func NewGoogleNewsAPI(apiKey string) *GoogleNewsAPI {
	return &GoogleNewsAPI{APIKey: apiKey}
}

func (api *GoogleNewsAPI) FetchTopHeadlinesNews(keyword string) ([]models.NewsArticle, error) {
	// hardcoding 'technology' because that's our main interest
	// we're using the 'top-headlines' path instead of 'everything' because it allows us to query further by country, category, etc.
	// news are sort by 'earliest date' from the api using above path
	// // https://newsapi.org/docs/endpoints/everything
	topHeadlinesUrl := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=us&category=technology&q=%s&apiKey=%s", keyword, api.APIKey)
	resp, err := utils.MakeSecureHTTPRequest(http.MethodGet, topHeadlinesUrl, nil)
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
		Articles []models.NewsArticle `json:"articles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding json: %s", err)
		return nil, err
	}

	log.Printf("Article Response: %s", result.Articles)

	return result.Articles, nil
}

// FetchEverythingHacking queries the News API 'everything' endpoint searching for hacking keywords
func (api *GoogleNewsAPI) FetchEverythingHacking() ([]models.NewsArticle, error) {
	return api.fetchEverythingNews(hackingQuery)
}

// fetchEverythingNews is the signature function used to hit the News API 'everything' endpoint. It expects
// an argument with the query keywords to use.
func (api *GoogleNewsAPI) fetchEverythingNews(query string) ([]models.NewsArticle, error) {
	// https://newsapi.org/docs/endpoints/everything
	encodedQuery := url.QueryEscape(query)
	// Check the length of the encoded query, max supported by api is 500 chars
	if len(encodedQuery) > 500 {
		return nil, fmt.Errorf("encoded query exceeds the maximum length of 500 characters")
	}

	// set date range
	//fromDate := time.Now().AddDate(0, 0, -15)
	//toDate := time.Now().Date()

	everythingUrl := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&apiKey=%s", encodedQuery, api.APIKey)
	log.Printf("Request everything url: %s", everythingUrl)
	resp, err := utils.MakeSecureHTTPRequest(http.MethodGet, everythingUrl, nil)
	if err != nil {
		log.Printf("Error making http get request: %s", err)
		return nil, err
	}
	// close connection
	defer func() {
		if err = resp.Body.Close(); err != nil {
			// Handle the error if needed, for example, log it
			log.Printf("failed to close response body: %v", err)
		}
	}()

	var result struct {
		Articles []models.NewsArticle `json:"articles"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Printf("Error decoding json: %s", err)
		return nil, err
	}

	log.Printf("Article Response: %s", result.Articles)

	return result.Articles, nil
}
