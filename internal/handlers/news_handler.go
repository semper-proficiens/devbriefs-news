package handlers

import (
	"devbriefs-news/internal/service"
	"encoding/json"
	"net/http"
)

type NewsHandler struct {
	newsService service.NewsServiceInterface
}

func NewNewsHandler(newsService service.NewsServiceInterface) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}

func (h *NewsHandler) GetEverythingHackingNews(w http.ResponseWriter, r *http.Request) {
	news, err := h.newsService.FetchEverythingHackingNews()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
