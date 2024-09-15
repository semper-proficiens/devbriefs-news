package handlers

import (
	"devbriefs-news/models"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockGoogleNewsAPI is a mock implementation of the NewsAPI interface.
type MockGoogleNewsAPI struct {
	FetchEverythingHackingFunc func() ([]models.NewsArticle, error)
}

func (m *MockGoogleNewsAPI) FetchEverythingHacking() ([]models.NewsArticle, error) {
	return m.FetchEverythingHackingFunc()
}

func TestGetEveryHackingNews(t *testing.T) {
	tests := []struct {
		name           string
		mockNews       []models.NewsArticle
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "FetchEverythingHacking returns error",
			mockNews:       nil,
			mockError:      errors.New("fetch error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "fetch error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAPI := &MockGoogleNewsAPI{
				FetchEverythingHackingFunc: func() ([]models.NewsArticle, error) {
					return tt.mockNews, tt.mockError
				},
			}

			req, err := http.NewRequest("GET", "/api/everything-hacking-news", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				GetEveryHackingNews(w, mockAPI)
			})

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if body := rr.Body.String(); body != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", body, tt.expectedBody)
			}
		})
	}
}
