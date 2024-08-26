package service

import (
	"news-fetcher/internal/api"
	"news-fetcher/internal/config"
	"news-fetcher/internal/models"
)

type NewsService struct {
	newsAPI *api.GoogleNewsAPI
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
