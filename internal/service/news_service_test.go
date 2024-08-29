package service

import (
	"devbriefs-news/internal/config"
	"devbriefs-news/internal/models"
	"errors"
	"testing"
)

// MockGoogleNewsAPI is a mock implementation of the GoogleNewsAPI interface.
type MockGoogleNewsAPI struct {
	FetchEverythingHackingFunc func() ([]models.NewsArticle, error)
}

func (m *MockGoogleNewsAPI) FetchEverythingHacking() ([]models.NewsArticle, error) {
	return m.FetchEverythingHackingFunc()
}

// TestNewNewsService tests the NewNewsService function.
func TestNewNewsService(t *testing.T) {
	cfg := &config.Config{
		GoogleAPIKey: "test-api-key",
	}
	service := NewNewsService(cfg)

	if service.newsAPI == nil {
		t.Error("Expected newsAPI to be non-nil")
	}
}

// TestFetchEverythingHackingNews_Success tests the FetchEverythingHackingNews method for a successful response.
func TestFetchEverythingHackingNews_Success(t *testing.T) {
	mockAPI := &MockGoogleNewsAPI{
		FetchEverythingHackingFunc: func() ([]models.NewsArticle, error) {
			return []models.NewsArticle{
				{Title: "Article 1", URL: "http://example.com/1"},
				{Title: "Article 2", URL: "http://example.com/2"},
			}, nil
		},
	}
	service := &NewsService{newsAPI: mockAPI}

	articles, err := service.FetchEverythingHackingNews()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(articles) != 2 {
		t.Fatalf("Expected 2 articles, got %d", len(articles))
	}

	expectedTitles := []string{"Article 1", "Article 2"}
	for i, article := range articles {
		if article.Title != expectedTitles[i] {
			t.Errorf("Expected title %s, got %s", expectedTitles[i], article.Title)
		}
	}
}

// TestFetchEverythingHackingNews_Error tests the FetchEverythingHackingNews method for an error response.
func TestFetchEverythingHackingNews_Error(t *testing.T) {
	mockAPI := &MockGoogleNewsAPI{
		FetchEverythingHackingFunc: func() ([]models.NewsArticle, error) {
			return nil, errors.New("test error")
		},
	}
	service := &NewsService{newsAPI: mockAPI}

	_, err := service.FetchEverythingHackingNews()
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	expectedError := "test error"
	if err.Error() != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, err.Error())
	}
}
