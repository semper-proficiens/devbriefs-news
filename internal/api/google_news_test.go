package api

import (
	"devbriefs-news/internal/models"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Mock function to replace utils.MakeSecureGetHTTPRequest
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

func TestNewGoogleNewsAPI(t *testing.T) {
	apiKey := "test-api-key"
	api := NewGoogleNewsAPI(apiKey)

	if api.APIKey != apiKey {
		t.Errorf("Expected APIKey to be %s, got %s", apiKey, api.APIKey)
	}

	if api.GetHTTPRequest == nil {
		t.Error("Expected GetHTTPRequest to be non-nil")
	}

	if api.URLParse == nil {
		t.Error("Expected URLParse to be non-nil")
	}
}

func TestFetchEverythingHacking(t *testing.T) {
	api := &GoogleNewsAPI{
		APIKey:         "test-api-key",
		GetHTTPRequest: mockGetHTTPRequest,
		URLParse:       url.Parse,
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

// Additional tests for error handling and edge cases

// Test FetchEverythingHacking with encoded query length error
func TestFetchEverythingHacking_EncodedQueryLengthError(t *testing.T) {
	api := &GoogleNewsAPI{
		APIKey:         "test-api-key",
		GetHTTPRequest: http.Get,
		URLParse:       url.Parse,
	}

	longQuery := string(make([]byte, 501))
	_, err := api.fetchEverythingNews(longQuery)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "encoded query exceeds the maximum length of 500 characters"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, err.Error())
	}
}

// Test FetchEverythingHacking with base URL parsing error
func TestFetchEverythingHacking_BaseURLParsingError(t *testing.T) {
	api := &GoogleNewsAPI{
		APIKey:         "test-api-key",
		GetHTTPRequest: mockGetHTTPRequest,
		URLParse: func(rawURL string) (*url.URL, error) {
			return nil, errors.New("mock parse error")
		},
	}

	_, err := api.fetchEverythingNews("test")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "failed to parse base URL: mock parse error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, err.Error())
	}
}

// Test FetchEverythingHacking with HTTP request error
func TestFetchEverythingHacking_HTTPError(t *testing.T) {
	api := &GoogleNewsAPI{
		APIKey: "test-api-key",
		GetHTTPRequest: func(url string) (*http.Response, error) {
			return nil, errors.New("mock error")
		},
		URLParse: url.Parse,
	}

	_, err := api.FetchEverythingHacking()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "mock error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, err.Error())
	}
}

// Test FetchEverythingHacking with JSON decoding error
func TestFetchEverythingHacking_JSONError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`invalid json`)); err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
	defer ts.Close()

	api := &GoogleNewsAPI{
		APIKey: "test-api-key",
		GetHTTPRequest: func(url string) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`invalid json`)),
				Header:     make(http.Header),
			}, nil
		},
		URLParse: url.Parse,
	}

	_, err := api.FetchEverythingHacking()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "invalid character 'i' looking for beginning of value"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, err.Error())
	}
}

// Test removeDuplicateTitles with similarity threshold
func TestRemoveDuplicateTitles(t *testing.T) {
	articles := []models.NewsArticle{
		{Title: "Article 2 Hack Something"},
		{Title: "Article 2 Hack Something Repeated Once"},
		{Title: "Article 2 Hack Something Repeated Twice"},
		{Title: "Something totally different"},
	}

	uniqueArticles := removeDuplicateTitles(articles, 0.6)
	if len(uniqueArticles) != 2 {
		t.Fatalf("Expected 2 unique articles, got %d", len(uniqueArticles))
	}
}

func BenchmarkFetchEverythingHacking(b *testing.B) {
	api := &GoogleNewsAPI{
		APIKey:         "test-api-key",
		GetHTTPRequest: mockGetHTTPRequest,
		URLParse:       url.Parse,
	}

	for i := 0; i < b.N; i++ {
		_, err := api.FetchEverythingHacking()
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
	}
}
