package handlers

import (
	"devbriefs-news/internal/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockNewsService is a mock implementation of the NewsServiceInterface.
type MockNewsService struct {
	FetchTopHeadlinesNewsFunc      func(keyword string) ([]models.NewsArticle, error)
	FetchEverythingHackingNewsFunc func() ([]models.NewsArticle, error)
}

func (m *MockNewsService) FetchTopHeadlinesNews(keyword string) ([]models.NewsArticle, error) {
	return m.FetchTopHeadlinesNewsFunc(keyword)
}

func (m *MockNewsService) FetchEverythingHackingNews() ([]models.NewsArticle, error) {
	return m.FetchEverythingHackingNewsFunc()
}

// TestGetTopHeadlinesNews_Success tests the GetTopHeadlinesNews method for a successful response.
func TestGetTopHeadlinesNews_Success(t *testing.T) {
	mockService := &MockNewsService{
		FetchTopHeadlinesNewsFunc: func(keyword string) ([]models.NewsArticle, error) {
			return []models.NewsArticle{
				{Title: "Article 1", URL: "http://example.com/1"},
				{Title: "Article 2", URL: "http://example.com/2"},
			}, nil
		},
	}
	handler := NewNewsHandler(mockService)

	req, err := http.NewRequest("GET", "/top-headlines?keyword=test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler.GetTopHeadlinesNews(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var actual []models.NewsArticle
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	expected := []models.NewsArticle{
		{Title: "Article 1", URL: "http://example.com/1"},
		{Title: "Article 2", URL: "http://example.com/2"},
	}

	if !equalArticles(actual, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

// TestGetTopHeadlinesNews_Error tests the GetTopHeadlinesNews method for an error response.
func TestGetTopHeadlinesNews_Error(t *testing.T) {
	mockService := &MockNewsService{
		FetchTopHeadlinesNewsFunc: func(keyword string) ([]models.NewsArticle, error) {
			return nil, errors.New("test error")
		},
	}
	handler := NewNewsHandler(mockService)

	req, err := http.NewRequest("GET", "/top-headlines?keyword=test", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler.GetTopHeadlinesNews(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := "test error\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestGetEverythingHackingNews_Success tests the GetEverythingHackingNews method for a successful response.
func TestGetEverythingHackingNews_Success(t *testing.T) {
	mockService := &MockNewsService{
		FetchEverythingHackingNewsFunc: func() ([]models.NewsArticle, error) {
			return []models.NewsArticle{
				{Title: "Article 1", URL: "http://example.com/1"},
				{Title: "Article 2", URL: "http://example.com/2"},
			}, nil
		},
	}
	handler := NewNewsHandler(mockService)

	req, err := http.NewRequest("GET", "/everything-hacking", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler.GetEverythingHackingNews(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var actual []models.NewsArticle
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	expected := []models.NewsArticle{
		{Title: "Article 1", URL: "http://example.com/1"},
		{Title: "Article 2", URL: "http://example.com/2"},
	}

	if !equalArticles(actual, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

// TestGetEverythingHackingNews_Error tests the GetEverythingHackingNews method for an error response.
func TestGetEverythingHackingNews_Error(t *testing.T) {
	mockService := &MockNewsService{
		FetchEverythingHackingNewsFunc: func() ([]models.NewsArticle, error) {
			return nil, errors.New("test error")
		},
	}
	handler := NewNewsHandler(mockService)

	req, err := http.NewRequest("GET", "/everything-hacking", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler.GetEverythingHackingNews(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}

	expected := "test error\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// equalArticles compares two slices of NewsArticle for equality.
func equalArticles(a, b []models.NewsArticle) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
