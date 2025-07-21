package themoviedb

import (
	"MediaTools/pkg/limiter"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"time"
)

type TMDB struct {
	baseURL string
	imgURl  string
	apiKey  string
	client  *http.Client
	limiter *limiter.Limiter
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
		baseURL: "https://api.themoviedb.org",
		client:  &http.Client{},
	}
	for _, opt := range opts {
		opt(config)
	}
	return &TMDB{
		baseURL: config.baseURL,
		imgURl:  "https://image.tmdb.org",
		apiKey:  apiKey,
		client:  config.client,
		limiter: limiter.NewLimiter(time.Second, 20),
	}
}

func (tmdb *TMDB) DoRequest(method string, path string, query url.Values, body io.Reader, resp any) error {
	query.Set("api_key", tmdb.apiKey)
	url := tmdb.baseURL + "/3" + path + "?" + query.Encode()
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
	if res.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("request failed, status code: %d, message: %s", res.StatusCode, errResp.StatusMessage)
		}
		return fmt.Errorf("decode error response failed: %w", err)
	}
	if err := json.NewDecoder(res.Body).Decode(resp); err != nil {
		return fmt.Errorf("decode response failed: %w", err)
	}
	return nil
}

func (tmdb *TMDB) DownloadImage(path string) (image.Image, error) {
	url := tmdb.imgURl + "/t/p/original" + path
	resp, err := tmdb.client.Get(url)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("下载图片「%s」失败", url))
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, NewTMDBError(nil, fmt.Sprintf("下载图片「%s」失败，HTTP code: %d", url, resp.StatusCode))
	}
	defer resp.Body.Close()
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("解码图片「%s」失败", url))
	}
	return img, nil
}
