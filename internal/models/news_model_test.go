package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestNewsArticleJSONMarshaling tests the JSON marshaling of the NewsArticle struct.
func TestNewsArticleJSONMarshaling(t *testing.T) {
	article := NewsArticle{
		Title:       "Test Title",
		URL:         "http://example.com",
		Description: "Test Description",
		Source:      NewsSource{ID: "test-id", Name: "Test Source"},
		PublishedAt: "2024-08-28T19:24:38Z",
	}

	data, err := json.Marshal(article)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedJSON := `{"title":"Test Title","url":"http://example.com","description":"Test Description","source":{"id":"test-id","name":"Test Source"},"publishedAt":"2024-08-28T19:24:38Z"}`
	if string(data) != expectedJSON {
		t.Errorf("Expected JSON %s, got %s", expectedJSON, string(data))
	}
}

// TestNewsArticleJSONUnmarshal tests the JSON unmarshal of the NewsArticle struct.
func TestNewsArticleJSONUnmarshal(t *testing.T) {
	jsonData := `{"title":"Test Title","url":"http://example.com","description":"Test Description","source":{"id":"test-id","name":"Test Source"},"publishedAt":"2024-08-28T19:24:38Z"}`

	var article NewsArticle
	err := json.Unmarshal([]byte(jsonData), &article)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedArticle := NewsArticle{
		Title:       "Test Title",
		URL:         "http://example.com",
		Description: "Test Description",
		Source:      NewsSource{ID: "test-id", Name: "Test Source"},
		PublishedAt: "2024-08-28T19:24:38Z",
	}

	if !reflect.DeepEqual(article, expectedArticle) {
		t.Errorf("Expected article %v, got %v", expectedArticle, article)
	}
}

// TestNewsSourceJSONMarshal tests the JSON marshaling of the NewsSource struct.
func TestNewsSourceJSONMarshal(t *testing.T) {
	source := NewsSource{
		ID:   "test-id",
		Name: "Test Source",
	}

	data, err := json.Marshal(source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedJSON := `{"id":"test-id","name":"Test Source"}`
	if string(data) != expectedJSON {
		t.Errorf("Expected JSON %s, got %s", expectedJSON, string(data))
	}
}

// TestNewsSourceJSONUnmarshal tests the JSON unmarshal of the NewsSource struct.
func TestNewsSourceJSONUnmarshal(t *testing.T) {
	jsonData := `{"id":"test-id","name":"Test Source"}`

	var source NewsSource
	err := json.Unmarshal([]byte(jsonData), &source)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedSource := NewsSource{
		ID:   "test-id",
		Name: "Test Source",
	}

	if !reflect.DeepEqual(source, expectedSource) {
		t.Errorf("Expected source %v, got %v", expectedSource, source)
	}
}
