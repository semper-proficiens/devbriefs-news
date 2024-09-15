package api

import (
	"context"
	"crypto/md5"
	"devbriefs-news/models"
	"devbriefs-news/services"
	"fmt"
	"github.com/semper-proficiens/go-utils/web/securehttp"
	"time"
)

// googleNewsTimeout is the time we want to make our go routines wait for a Google News API response in milliseconds
const googleNewsTimeout = 500

type NewsAPIResponse struct {
	articles []models.NewsArticle
	err      error
}

// NewsAPI defines the interface for fetching news articles.
type NewsAPI interface {
	FetchEverythingHacking(ctx context.Context) (map[string]any, error)
}

type GoogleNewsAPI struct {
	APIKey     string
	HTTPClient securehttp.CustomHTTPClientInterface
}

func NewGoogleNewsAPI(apiKey string, sc *securehttp.CustomHTTPClient) (*GoogleNewsAPI, error) {
	return &GoogleNewsAPI{
		APIKey:     apiKey,
		HTTPClient: sc,
	}, nil
}

// FetchEverythingHacking is an API method that calls the "FetchEverythingLogic" service logic for "hacking"
func (api *GoogleNewsAPI) FetchEverythingHacking(ctx context.Context) (map[string]any, error) {
	// we'll cancel this operation if it exceeds this time
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*googleNewsTimeout)
	defer cancel()

	hackingChan := make(chan NewsAPIResponse)
	data := make(map[string]any)

	go func() {
		a, err := services.FetchEverythingNews(ctx, "hacking", api.APIKey, api.HTTPClient)
		hackingChan <- NewsAPIResponse{
			articles: a,
			err:      err,
		}
	}()

	// blocking until go routine context expires or we get a response from the api
	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("FetchEverythingHacking timed out after %d milliseconds", googleNewsTimeout)
		case apiResponse := <-hackingChan:
			for _, article := range apiResponse.articles {
				hashedTitle := fmt.Sprintf("%x", md5.Sum([]byte(article.Title)))
				data[hashedTitle] = article
			}
			return data, apiResponse.err
		}
	}
}
