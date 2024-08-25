package models

// NewsArticleEverything represents a single news article fetched from the Google News API from 'everything' endpoint
type NewsArticleEverything struct {
	Title       string `json:"title"`       // The title of the news article
	URL         string `json:"url"`         // The URL to the full news article
	Description string `json:"description"` // A brief description of the news article
	Source      string `json:"source"`      // The source of the news article
	PublishedAt string `json:"publishedAt"` // The publication date of the news article
}

// NewsArticleTopHeadlines represents a single news article fetched from the Google News API from 'top-headlines' endpoint
type NewsArticleTopHeadlines struct {
	Title       string             `json:"title"`       // The title of the news article
	URL         string             `json:"url"`         // The URL to the full news article
	Description string             `json:"description"` // A brief description of the news article
	Source      TopHeadlinesSource `json:"source"`      // The source of the news article
	PublishedAt string             `json:"publishedAt"` // The publication date of the news article
}

// TopHeadlinesSource represents the source of a news article.
type TopHeadlinesSource struct {
	ID   string `json:"id"`   // The ID of the news source
	Name string `json:"name"` // The name of the news source
}
