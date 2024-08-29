package service

import (
	"devbriefs-news/internal/api"
	"devbriefs-news/internal/config"
	"devbriefs-news/internal/models"
)

// NewsServiceInterface defines the methods that the NewsService should implement.
type NewsServiceInterface interface {
	FetchTopHeadlinesNews(keyword string) ([]models.NewsArticle, error)
	FetchEverythingHackingNews() ([]models.NewsArticle, error)
}

type NewsService struct {
	newsAPI api.GoogleNewsAPIInterface
}

func NewNewsService(cfg *config.Config) *NewsService {
	return &NewsService{
		newsAPI: api.NewGoogleNewsAPI(cfg.GoogleAPIKey),
	}
}

func (s *NewsService) FetchTopHeadlinesNews(keyword string) ([]models.NewsArticle, error) {
	return s.newsAPI.FetchTopHeadlinesNews(keyword)
}

func (s *NewsService) FetchEverythingHackingNews() ([]models.NewsArticle, error) {
	return s.newsAPI.FetchEverythingHacking()
}
