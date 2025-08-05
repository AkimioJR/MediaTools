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

	"github.com/allegro/bigcache"
)

type TMDB struct {
	apiURL        string
	imgURL        string
	apiKey        string
	language      string
	imageLanguage string
	client        *http.Client
	limiter       *limiter.Limiter
	cache         *bigcache.BigCache
}

type tmdbConfig struct {
	apiURL   string
	imgURL   string
	language string
	client   *http.Client
	limiter  *limiter.Limiter
}
type TMDBOptions func(c *tmdbConfig)

func CustomAPIURL(apiURL string) TMDBOptions {
	return func(c *tmdbConfig) {
		c.apiURL = apiURL
	}
}

func CustomImageURL(imgURL string) TMDBOptions {
	return func(c *tmdbConfig) {
		c.imgURL = imgURL
	}
}

func CustomHTTPClient(client *http.Client) TMDBOptions {
	return func(c *tmdbConfig) {
		c.client = client
	}
}

func CustomLanguage(language string) TMDBOptions {
	return func(c *tmdbConfig) {
		c.language = language
	}
}

func CustomLimiter(d time.Duration, maxCount uint64) TMDBOptions {
	return func(c *tmdbConfig) {
		c.limiter = limiter.NewLimiter(d, maxCount)
	}
}

func NewTMDB(apiKey string, opts ...TMDBOptions) *TMDB {
	config := &tmdbConfig{
		apiURL: "https://api.themoviedb.org",
		client: &http.Client{},
	}

	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))

	for _, opt := range opts {
		opt(config)
	}
	client := TMDB{
		apiURL:        config.apiURL,
		imgURL:        "https://image.tmdb.org",
		apiKey:        apiKey,
		language:      "zh-CN",
		imageLanguage: "zh",
		client:        config.client,
		limiter:       limiter.NewLimiter(time.Second, 20),
		cache:         cache,
	}
	return &client
}

func (tmdb *TMDB) DoRequest(method string, path string, query url.Values, body io.Reader, resp any) error {
	var (
		data []byte
		err  error
	)
	cacheKey := method + "|" + path + "|" + query.Encode()
	if method == http.MethodGet && body == nil {
		if data, err = tmdb.cache.Get(cacheKey); err == nil {
			return json.Unmarshal(data, resp)
		}
	}

	tmdb.limiter.Acquire()
	query.Set("api_key", tmdb.apiKey)
	url := tmdb.apiURL + "/3" + path + "?" + query.Encode()
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
	data, err = io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err)
	}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return fmt.Errorf("unmarshal response failed: %w", err)
	}

	// 写入缓存
	if method == http.MethodGet && body == nil {
		if data, err = json.Marshal(resp); err == nil {
			_ = tmdb.cache.Set(cacheKey, data)
		}
	}
	return nil
}

func (tmdb *TMDB) GetImageURL(path string) string {
	return tmdb.imgURL + "/t/p/original" + path
}

func (tmdb *TMDB) DownloadImage(path string) (image.Image, error) {
	url := tmdb.GetImageURL(path)
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
