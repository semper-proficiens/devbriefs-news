package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func mockGetHTTPRequest(url string) (*http.Response, error) {
	var mockResponse string

	if strings.Contains(url, "top-headlines") {
		mockResponse = `{"articles": [{"title": "Test Article", "url": "https://example.com", "publishedAt": "2023-10-01T00:00:00Z"}]}`
	} else if strings.Contains(url, "everything") {
		mockResponse = `{"articles": [{"title": "Hacking News", "url": "https://example.com", "publishedAt": "2023-10-01T00:00:00Z"}]}`
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(mockResponse)),
		Header:     make(http.Header),
	}, nil
}

func TestFetchTopHeadlinesNews(t *testing.T) {
	api := &GoogleNewsAPI{
		APIKey:         "test-api-key",
		GetHTTPRequest: mockGetHTTPRequest,
	}

	articles, err := api.FetchTopHeadlinesNews("test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(articles) != 1 {
		t.Fatalf("Expected 1 article, got %d", len(articles))
	}

	if articles[0].Title != "Test Article" {
		t.Errorf("Expected article title to be 'Test Article', got %s", articles[0].Title)
	}
}

func TestFetchEverythingHacking(t *testing.T) {
	api := &GoogleNewsAPI{
		APIKey:         "test-api-key",
		GetHTTPRequest: mockGetHTTPRequest,
	}

	articles, err := api.FetchEverythingHacking()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(articles) != 1 {
		t.Fatalf("Expected 1 article, got %d", len(articles))
	}

	if articles[0].Title != "Hacking News" {
		t.Errorf("Expected article title to be 'Hacking News', got %s", articles[0].Title)
	}
}

func BenchmarkFetchTopHeadlinesNews(b *testing.B) {
	api := &GoogleNewsAPI{
		APIKey:         "test-api-key",
		GetHTTPRequest: mockGetHTTPRequest,
	}

	for i := 0; i < b.N; i++ {
		_, err := api.FetchTopHeadlinesNews("test")
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
	}
}

func BenchmarkFetchEverythingHacking(b *testing.B) {
	api := &GoogleNewsAPI{
		APIKey:         "test-api-key",
		GetHTTPRequest: mockGetHTTPRequest,
	}

	for i := 0; i < b.N; i++ {
		_, err := api.FetchEverythingHacking()
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
	}
}
