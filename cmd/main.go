package main

import (
	"devbriefs-news/internal/config"
	"devbriefs-news/internal/handlers"
	"devbriefs-news/internal/service"
	"devbriefs-news/internal/utils"
	"log"

	"github.com/gin-gonic/gin"
)

const cloudFlareAPI = "https://api.cloudflare.com/client/v4"

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Fetch Cloudflare IP ranges once during startup
	trustedProxies, err := utils.FetchCloudflareIPv4Ranges(cloudFlareAPI)
	if err != nil {
		log.Fatalf("Failed to fetch Cloudflare IP ranges: %v", err)
	}

	// Initialize news service
	newsService := service.NewNewsService(cfg)

	// Initialize handlers
	newsHandler := handlers.NewNewsHandler(newsService)

	// Set up Gin router for production
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/api/everything-hacking-news", func(c *gin.Context) {
		newsHandler.GetEverythingHackingNews(c.Writer, c.Request)
	})

	// Set trusted proxies to Cloudflare IP ranges
	if err = r.SetTrustedProxies(trustedProxies); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Start server using Gin's built-in method
	log.Println("Starting server on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
