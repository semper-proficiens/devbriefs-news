package handlers

import (
	"devbriefs-news/internal/api"
	"encoding/json"
	"net/http"
)

func GetEveryHackingNews(w http.ResponseWriter, api api.NewsAPI) {
	news, err := api.FetchEverythingHacking()
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
