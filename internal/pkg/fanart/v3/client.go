package fanart

import (
	"MediaTools/pkg/limiter"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/allegro/bigcache"
)

type FanartClient struct {
	api     string
	apiKey  string
	client  *http.Client
	limiter *limiter.Limiter
	cache   *bigcache.BigCache
}

func NewClient(apiKey string) (*FanartClient, error) {
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		return nil, fmt.Errorf("create cache for FanartClient failed: %w", err)
	}

	client := FanartClient{
		api:     "https://webservice.fanart.tv",
		apiKey:  apiKey,
		client:  &http.Client{},
		limiter: limiter.NewLimiter(1*time.Second, 20), // 每秒最多20次请求
		cache:   cache,
	}
	return &client, nil
}

func (client *FanartClient) DoRequest(method string, path string, query url.Values, body io.Reader, resp any) error {
	var (
		data []byte
		err  error
	)
	cacheKey := method + "|" + path + "|" + query.Encode()
	if method == http.MethodGet && body == nil {
		if data, err = client.cache.Get(cacheKey); err == nil {
			return json.Unmarshal(data, resp)
		}
	}

	client.limiter.Acquire()
	query.Set("api_key", client.apiKey)
	url := client.api + "/3" + path + "?" + query.Encode()
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	res, err := client.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request failed: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("request failed, status code: %d, message: %s", res.StatusCode, errResp.ErrorMessage)
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
			_ = client.cache.Set(cacheKey, data)
		}
	}
	return nil
}
