package api

import (
	"devbriefs-news/internal/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

// MockHTTPClient is a mock implementation of the CustomHTTPClientInterface.
type MockHTTPClient struct {
	GetFunc func(url string) (*http.Response, error)
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}

func TestFetchEverythingHacking(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockGetFunc    func(url string) (*http.Response, error)
		mockURLParse   func(rawURL string) (*url.URL, error)
		expectedResult []models.NewsArticle
		expectedError  error
	}{
		{
			name:  "Successful fetch and parse",
			query: hackingQuery,
			mockGetFunc: func(url string) (*http.Response, error) {
				articles := []models.NewsArticle{
					{Title: "Test News 1"},
					{Title: "Test News 2"},
				}
				body, _ := json.Marshal(map[string]interface{}{"articles": articles})
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(string(body))),
				}, nil
			},
			mockURLParse: func(rawURL string) (*url.URL, error) {
				return url.Parse(rawURL)
			},
			expectedResult: []models.NewsArticle{
				{Title: "Test News 1"},
			},
			expectedError: nil,
		},
		{
			name:  "URL parsing fails",
			query: hackingQuery,
			mockGetFunc: func(url string) (*http.Response, error) {
				return nil, nil
			},
			mockURLParse: func(rawURL string) (*url.URL, error) {
				return nil, errors.New("url parse error")
			},
			expectedResult: nil,
			expectedError:  errors.New("failed to parse base URL: url parse error"),
		},
		{
			name:  "HTTP request fails",
			query: hackingQuery,
			mockGetFunc: func(url string) (*http.Response, error) {
				return nil, errors.New("http request error")
			},
			mockURLParse: func(rawURL string) (*url.URL, error) {
				return url.Parse(rawURL)
			},
			expectedResult: nil,
			expectedError:  errors.New("http request error"),
		},
		{
			name:  "JSON unmarshal fails",
			query: hackingQuery,
			mockGetFunc: func(url string) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader("invalid json")),
				}, nil
			},
			mockURLParse: func(rawURL string) (*url.URL, error) {
				return url.Parse(rawURL)
			},
			expectedResult: nil,
			expectedError:  errors.New("invalid character 'i' looking for beginning of value"),
		},
		{
			name:  "Encoded query exceeds max length",
			query: strings.Repeat("a", 501),
			mockGetFunc: func(url string) (*http.Response, error) {
				return nil, nil
			},
			mockURLParse: func(rawURL string) (*url.URL, error) {
				return url.Parse(rawURL)
			},
			expectedResult: nil,
			expectedError:  errors.New("encoded query exceeds the maximum length of 500 characters"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHTTPClient := &MockHTTPClient{
				GetFunc: tt.mockGetFunc,
			}

			api := &GoogleNewsAPI{
				APIKey:     "test-api-key",
				HTTPClient: mockHTTPClient,
				URLParse:   tt.mockURLParse,
			}

			result, err := api.fetchEverythingNews(tt.query)

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("expected result %v, got %v", tt.expectedResult, result)
			}

			if err != nil && tt.expectedError != nil && err.Error() != tt.expectedError.Error() {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			} else if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}
