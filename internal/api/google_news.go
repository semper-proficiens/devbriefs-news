package api

import (
	"devbriefs-news/internal/models"
	"devbriefs-news/internal/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	hackingQuery = `
	"data breach" OR 
	"hacker" OR 
	"hackers" OR 
	"hacked" OR 
	"malware" OR 
	"exploited vulnerability"
	-"how to"
	-"your"
	-"you"
	-"my"
	`
	newsLanguage        = "en"
	newsExcludedDomains = `
	daemonology.net
	`
	//only root domains, not fqdn (e.g. talosintelligence.com vs blog.talosintelligence.com)
	newsDomains  = "thehackernews.com,hackread.com,talosintelligence.com,bleepingcomputer.com,cisa.gov,csoonline.com,threatpost.com,krebsonsecurity.com,wired.com,zdnet.com,virtualattacks.com"
	newsSortBy   = "publishedAt" // options: "relevancy" to q, "publishedAt" for newest (default)
	newsPageSize = "10"
)

type GoogleNewsAPI struct {
	APIKey         string
	GetHTTPRequest func(url string) (*http.Response, error)
}

func NewGoogleNewsAPI(apiKey string) *GoogleNewsAPI {
	return &GoogleNewsAPI{
		APIKey:         apiKey,
		GetHTTPRequest: utils.MakeSecureGetHTTPRequest,
	}
}

func (api *GoogleNewsAPI) FetchTopHeadlinesNews(keyword string) ([]models.NewsArticle, error) {
	// hardcoding 'technology' because that's our main interest
	// we're using the 'top-headlines' path instead of 'everything' because it allows us to query further by country, category, etc.
	// news are sort by 'earliest date' from the api using above path
	// Official doc https://newsapi.org/docs/endpoints/everything
	topHeadlinesUrl := fmt.Sprintf("https://newsapi.org/v2/top-headlines?country=us&category=technology&q=%s&apiKey=%s", keyword, api.APIKey)
	resp, err := api.GetHTTPRequest(topHeadlinesUrl)
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

	//log.Printf("Article Response: %s", result.Articles)

	return result.Articles, nil
}

// FetchEverythingHacking queries the News API 'everything' endpoint searching for hacking keywords
func (api *GoogleNewsAPI) FetchEverythingHacking() ([]models.NewsArticle, error) {
	return api.fetchEverythingNews(hackingQuery)
}

// fetchEverythingNews is the signature function used to hit the News API 'everything' endpoint. It expects
// an argument with the query keywords to use.
// Official doc https://newsapi.org/docs/endpoints/everything
func (api *GoogleNewsAPI) fetchEverythingNews(query string) ([]models.NewsArticle, error) {
	encodedQuery := url.QueryEscape(query)
	// Check the length of the encoded query, max supported by api is 500 chars
	if len(encodedQuery) > 500 {
		return nil, fmt.Errorf("encoded query exceeds the maximum length of 500 characters")
	}

	baseURL, err := url.Parse("https://newsapi.org/v2/everything")
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	// get news from 1 week ago
	fromDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	// query object
	params := url.Values{}
	params.Add("q", hackingQuery)
	params.Add("searchin", "title")
	params.Add("excludeDomains", newsExcludedDomains)
	params.Add("language", newsLanguage)
	params.Add("sortBy", newsSortBy)
	params.Add("domains", newsDomains)
	params.Add("pageSize", newsPageSize)
	params.Add("from", fromDate)
	params.Add("to", toDate)
	params.Add("apiKey", api.APIKey)

	// Add the query parameters to the URL
	baseURL.RawQuery = params.Encode()

	resp, err := api.GetHTTPRequest(baseURL.String())
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

	uniqueArticles := removeDuplicateTitles(result.Articles, 0.6)

	//for _, ua := range uniqueArticles {
	//	log.Printf("Title: %s, URL: %s, Date: %s", ua.Title, ua.URL, ua.PublishedAt)
	//}

	return uniqueArticles, nil
}

// removeDuplicateTitles removes duplicate titles based on a similarity threshold.
func removeDuplicateTitles(articles []models.NewsArticle, threshold float64) []models.NewsArticle {
	var uniqueArticles []models.NewsArticle
	for i, article := range articles {
		isDuplicate := false
		for j := 0; j < i; j++ {
			similarity := utils.CalculateSimilarity(article.Title, articles[j].Title)
			if similarity >= threshold {
				log.Printf("Excluded potential duplicate news \"%s\" and \"%s\" with similarity score %2f", article.Title, articles[j].Title, similarity)
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			uniqueArticles = append(uniqueArticles, article)
		}
	}
	return uniqueArticles
}
