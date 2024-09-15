package main

import (
	"context"
	"devbriefs-news/api"
	"devbriefs-news/datastore"
	"devbriefs-news/handlers"
	"devbriefs-news/models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/semper-proficiens/go-utils/system/config"
	utilTime "github.com/semper-proficiens/go-utils/system/time"
	"github.com/semper-proficiens/go-utils/web/jsonhandler"
	"github.com/semper-proficiens/go-utils/web/securehttp"
	"log"
	"time"
)

const cloudFlareAPI = "https://api.cloudflare.com/client/v4/ips"

func main() {
	// Load configuration
	envVars := config.LoadEnvVars()
	googleAPIKey := envVars["GOOGLE_NEWS_API_KEY"]

	// let's instantiate our custom secure client
	sc, err := securehttp.NewSecureHTTPClient()
	if err != nil {
		log.Fatalf("failed to create secure http client: %v", err)
	}

	// pass key and secure client to our google-news api
	googleNewAPI, err := api.NewGoogleNewsAPI(googleAPIKey, sc)
	if err != nil {
		log.Fatalf("failed to create google api: %v", err)
	}

	// Set up Gin router for production
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// start our main context
	ctx := context.Background()

	// init cache
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "192.168.0.229:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer func() {
		if err = redisClient.Close(); err != nil {
			log.Fatalf("failed to close redis client: %v", err)
		}
	}()

	redisCache := datastore.NewRedisCache(redisClient)

	//redisCache.Scan()

	//if err = redisCache.Set("key0", "value0"); err != nil {
	//	log.Fatalf("Failed to SET in Cache: %v", err)
	//}
	//log.Println("Set key and value in Cache")

	//val, err := redisCache.Get("57e631d5-af69-4856-a49d-1af566149f35")
	//if err != nil {
	//	log.Println("failed to GET from Cache:", err)
	//}
	//log.Println("Got value from Cache:", val)

	//if err = redisCache.Remove("key0"); err != nil {
	//	log.Fatalf("Failed to remove key0: %v", err)
	//}
	//log.Println("Removed key from cache")

	//if err = redisCache.DeleteAll(); err != nil {
	//	log.Fatalf("failed to delete all redis cache keys in the current DB: %v", err)
	//}
	//log.Println("deleted all keys in the current DB")

	// we want to run this every day at 6am EST
	waitTime, err := utilTime.TimeUntilNextRun("America/New_York", 06, 00)
	if err != nil {
		log.Println("failed to obtain a valid wait time:", err)
	}

	var news map[string]any
	go func() {
		for {
			log.Println("Sleeping for", waitTime)
			time.Sleep(waitTime)
			news, err = googleNewAPI.FetchEverythingHacking(ctx)
			if err != nil {
				log.Fatalf("failed to fetch hacking news in daily routine: %v", err)
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

			log.Println("News Articles were refreshed as part of daily routine")

			waitTime, err = utilTime.TimeUntilNextRun("America/New_York", 00, 20)
			if err != nil {
				log.Println("failed to obtain a valid wait time:", err)
			}
		}
	}()

	r.GET("/api/everything-hacking-news", func(c *gin.Context) {
		handlers.GetEveryHackingNews(ctx, c.Writer, googleNewAPI, redisCache)
	})

	// let's make sure we're always getting valid CloudFlare IPv4 addresses
	// to initiate our gin router allowed proxies
	resp, err := sc.Get(cloudFlareAPI)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}
	defer securehttp.ResponseBodyCloser(resp.Body)

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
