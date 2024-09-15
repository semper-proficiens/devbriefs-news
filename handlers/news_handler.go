package handlers

import (
	"context"
	"devbriefs-news/api"
	"devbriefs-news/datastore"
	"encoding/json"
	"log"
	"net/http"
)

func GetEveryHackingNews(ctx context.Context, w http.ResponseWriter, api api.NewsAPI, redisCache *datastore.RedisCache) {
	news, err := api.FetchEverythingHacking(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store news in Cache
	for k, v := range news {
		jsonValue, err := json.Marshal(v)
		if err != nil {
			log.Printf("failed to marshal news data with key (%s) and value(%v): %v", k, v, err)
			continue
		}
		if err = redisCache.Set(k, jsonValue); err != nil {
			log.Println("failed to store news data:", err)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
