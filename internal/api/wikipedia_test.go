package api

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"
)

// mockHTTPClient is a mock of the HTTP Client
type mockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}

// TestGetArticleViews tests the GetArticleViews function
func TestGetArticleViews(t *testing.T) {
    // Create a test server simulating the Wikipedia API
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintln(w, `{"items":[{"views": 123456}]}`)
    }))
    defer ts.Close()

    // Create a mock HTTP client for the Wikipedia API
    client := &mockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            return http.Get(ts.URL)
        },
    }

    wikipediaAPI := NewWikipediaAPI(client)

    views, err := wikipediaAPI.GetArticleViews(context.Background(), "Go_(programming_language)", 2023, 4)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }

    if views != 123456 {
        t.Errorf("Expected 123456 views, got %d", views)
    }
}
