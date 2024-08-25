package models

// NewsArticle represents a single news article fetched from the Google News API.
type NewsArticle struct {
	Title       string `json:"title"`       // The title of the news article
	URL         string `json:"url"`         // The URL to the full news article
	Description string `json:"description"` // A brief description of the news article
	Source      string `json:"source"`      // The source of the news article
	PublishedAt string `json:"publishedAt"` // The publication date of the news article
}
