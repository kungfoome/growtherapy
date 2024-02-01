package httpclient

import (
    "net/http"
)

// Client is an interface to make testing easier.
type Client interface {
    Do(req *http.Request) (*http.Response, error)
}

// NewClient returns a new HTTP client.
func NewClient() Client {
    return &http.Client{}
}

// Get makes a GET request and returns the response body.
func Get(client Client, url string) (*http.Response, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    return client.Do(req)
}
