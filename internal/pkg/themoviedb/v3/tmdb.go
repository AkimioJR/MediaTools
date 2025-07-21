package themoviedb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type TMDB struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type tmdbConfig struct {
	baseURL string
	client  *http.Client
}
type TMDBOptions func(c *tmdbConfig)

func CustomBaseURL(baseURL string) TMDBOptions {
	return func(c *tmdbConfig) {
		c.baseURL = baseURL
	}
}

func CustomHTTPClient(client *http.Client) TMDBOptions {
	return func(c *tmdbConfig) {
		c.client = client
	}
}

func NewTMDB(apiKey string, opts ...TMDBOptions) *TMDB {
	config := &tmdbConfig{
		baseURL: "https://api.themoviedb.org/3",
		client:  &http.Client{},
	}
	for _, opt := range opts {
		opt(config)
	}
	return &TMDB{
		baseURL: config.baseURL,
		apiKey:  apiKey,
		client:  config.client,
	}
}

func (tmdb *TMDB) NewRequest(method, path string, query url.Values, body io.Reader, resp any) error {
	query.Set("api_key", tmdb.apiKey)
	url := tmdb.baseURL + path + "?" + query.Encode()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := tmdb.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request failed: %w", err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return fmt.Errorf("decode response failed: %w", err)
	}
	return nil
}
