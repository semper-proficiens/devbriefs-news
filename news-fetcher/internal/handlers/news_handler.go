package handlers

import (
	"encoding/json"
	"net/http"
	"news-fetcher/internal/service"
)

type NewsHandler struct {
	newsService *service.NewsService
}

func NewNewsHandler(newsService *service.NewsService) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}

func (h *NewsHandler) GetNews(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	keyword := r.URL.Query().Get("keyword")

	news, err := h.newsService.FetchNews(category, keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(news)
}
