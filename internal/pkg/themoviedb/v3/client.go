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

type Client struct {
	apiURL        string
	imgURL        string
	apiKey        string
	language      string
	imageLanguage string
	client        *http.Client
	limiter       *limiter.Limiter
	cache         *bigcache.BigCache
}

type clientConfig struct {
	apiURL   string
	imgURL   string
	language string
	client   *http.Client
	limiter  *limiter.Limiter
}
type ClientOptions func(c *clientConfig)

func CustomAPIURL(apiURL string) ClientOptions {
	return func(c *clientConfig) {
		c.apiURL = apiURL
	}
}

func CustomImageURL(imgURL string) ClientOptions {
	return func(c *clientConfig) {
		c.imgURL = imgURL
	}
}

func CustomHTTPClient(client *http.Client) ClientOptions {
	return func(c *clientConfig) {
		c.client = client
	}
}

func CustomLanguage(language string) ClientOptions {
	return func(c *clientConfig) {
		c.language = language
	}
}

func CustomLimiter(d time.Duration, maxCount uint64) ClientOptions {
	return func(c *clientConfig) {
		c.limiter = limiter.NewLimiter(d, maxCount)
	}
}

func NewClient(apiKey string, opts ...ClientOptions) (*Client, error) {
	config := &clientConfig{
		apiURL: "https://api.themoviedb.org",
		client: &http.Client{},
	}

	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		return nil, fmt.Errorf("创建 TMDB 缓存失败: %w", err)
	}
	for _, opt := range opts {
		opt(config)
	}
	client := Client{
		apiURL:        config.apiURL,
		imgURL:        "https://image.tmdb.org",
		apiKey:        apiKey,
		language:      "zh-CN",
		imageLanguage: "zh",
		client:        config.client,
		limiter:       limiter.NewLimiter(time.Second, 20),
		cache:         cache,
	}
	return &client, nil
}

func (c *Client) DoRequest(method string, path string, query url.Values, body io.Reader, resp any) error {
	var (
		data []byte
		err  error
	)
	cacheKey := method + "|" + path + "|" + query.Encode()
	if method == http.MethodGet && body == nil {
		if data, err = c.cache.Get(cacheKey); err == nil {
			return json.Unmarshal(data, resp)
		}
	}

	c.limiter.Acquire()
	query.Set("api_key", c.apiKey)
	url := c.apiURL + "/3" + path + "?" + query.Encode()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
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
			_ = c.cache.Set(cacheKey, data)
		}
	}
	return nil
}

func (c *Client) GetImageURL(path string) string {
	return c.imgURL + "/t/p/original" + path
}

func (c *Client) DownloadImage(path string) (image.Image, error) {
	url := c.GetImageURL(path)
	resp, err := c.client.Get(url)
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
