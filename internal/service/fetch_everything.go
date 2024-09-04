package service

import (
	"context"
	"devbriefs-news/internal/models"
	"github.com/semper-proficiens/go-utils/nlp"
	"github.com/semper-proficiens/go-utils/web/jsonhandler"
	"github.com/semper-proficiens/go-utils/web/securehttp"
	"github.com/semper-proficiens/go-utils/web/urlcleaner"
	"log"
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
	newsLanguage = "en"
	//only root domains, not fqdn (e.g. talosintelligence.com vs blog.talosintelligence.com)
	newsDomains  = "thehackernews.com,hackread.com,talosintelligence.com,bleepingcomputer.com,cisa.gov,csoonline.com,threatpost.com,krebsonsecurity.com,wired.com,zdnet.com,virtualattacks.com"
	newsSortBy   = "publishedAt" // options: "relevancy" to q, "publishedAt" for newest (default)
	newsPageSize = "10"
)

// FetchEverythingNews is function used to hit the News API 'everything' endpoint. It expects
// a newsType argument that will be mapped to some query logic associated to that newsType.
//
// For now only "hacking" news logic exists, so "newsType" will only accept "hacking" as an argument, and will default to
// a hacking query.
//
// e.g. FetchEverythingNews(ctx, "hacking")
// Official doc https://newsapi.org/docs/endpoints/everything
func FetchEverythingNews(ctx context.Context, newsType string, apiKey string, client securehttp.CustomHTTPClientInterface) ([]models.NewsArticle, error) {

	// by default newsType will be a hacking query
	var query string
	switch newsType {
	case "hacking":
		query = hackingQuery
	default:
		query = hackingQuery
	}

	baseURL, err := urlcleaner.UrlParser(query, "https://newsapi.org/v2/everything", 500)
	if err != nil {
		return nil, err
	}

	// get news from 1 week ago
	fromDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	toDate := time.Now().Format("2006-01-02")

	// query object
	params := url.Values{}
	params.Add("q", hackingQuery)
	params.Add("searchin", "title")
	params.Add("language", newsLanguage)
	params.Add("sortBy", newsSortBy)
	params.Add("domains", newsDomains)
	params.Add("pageSize", newsPageSize)
	params.Add("from", fromDate)
	params.Add("to", toDate)
	params.Add("apiKey", apiKey)

	// Add the query parameters to the URL
	baseURL.RawQuery = params.Encode()

	resp, err := client.Get(baseURL.String())
	if err != nil {
		return nil, err
	}
	// close connection
	defer func() {
		if err = resp.Body.Close(); err != nil {
			// Handle the error if needed, for example, log it
			log.Printf("failed to close response body: %v", err)
		}
	}()

	// let's unmarshal that response from API
	var result struct {
		Articles []models.NewsArticle `json:"articles"`
	}
	err = jsonhandler.UnmarshalJSONResponse(resp, &result)
	if err != nil {
		return nil, err
	}

	uniqueArticles := nlp.RemoveDuplicates(result.Articles, 0.6, "Title")

	for _, ua := range uniqueArticles {
		log.Printf("Title: %s, URL: %s, Date: %s", ua.Title, ua.URL, ua.PublishedAt)
	}

	return uniqueArticles, nil
}
