package fanart

import "net/http"

type options struct {
	apiURL string // Fanart API URL
	client *http.Client
}

type Options func(opt *options)

func CustomAPIURL(apiURL string) Options {
	return func(opt *options) {
		opt.apiURL = apiURL
	}
}

func CustomHTTPClient(client *http.Client) Options {
	return func(opt *options) {
		opt.client = client
	}
}
