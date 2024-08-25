package main

import (
	"log"
	"net/http"
	"news-fetcher/internal/config"
	"news-fetcher/internal/handlers"
	"news-fetcher/internal/service"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize news service
	newsService := service.NewNewsService(cfg)

	// Initialize handlers
	newsHandler := handlers.NewNewsHandler(newsService)

	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/api/news", newsHandler.GetNews).Methods("GET")

	// Start server
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
