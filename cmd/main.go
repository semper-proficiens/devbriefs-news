package main

import (
	"context"
	"devbriefs-news/internal/api"
	"devbriefs-news/internal/config"
	"devbriefs-news/internal/handlers"
	"devbriefs-news/internal/models"
	"github.com/semper-proficiens/go-utils/web/jsonhandler"
	"github.com/semper-proficiens/go-utils/web/securehttp"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

const cloudFlareAPI = "https://api.cloudflare.com/client/v4/ips"

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// define context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// let's instantiate our custom secure client
	sc, err := securehttp.NewSecureHTTPClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create secure http client: %v", err)
	}

	// pass key and secure client to our google-news api
	googleNewAPI, err := api.NewGoogleNewsAPI(cfg.GoogleAPIKey, sc)
	if err != nil {
		log.Fatalf("Failed to create google api: %v", err)
	}

	// Set up Gin router for production
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/api/everything-hacking-news", func(c *gin.Context) {
		handlers.GetEveryHackingNews(c.Writer, googleNewAPI)
	})

	// let's make sure we're always getting valid CloudFlare IPv4 addresses
	// to initiate our gin router allowed proxies
	resp, err := sc.Get(cloudFlareAPI)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			// Handle the error if needed, for example, log it
			log.Printf("failed to close response body: %v", err)
		}
	}()
	var ipRanges models.CloudflareIPRanges
	if err = jsonhandler.UnmarshalJSONResponse(resp, &ipRanges); err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Set trusted proxies to Cloudflare IP ranges
	if err = r.SetTrustedProxies(ipRanges.Result.IPv4CIDRs); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Start server using Gin's built-in method
	log.Println("Starting server on :8080")
	if err = r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
