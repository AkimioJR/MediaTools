package tmdb_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/outbound"
	"MediaTools/internal/pkg/themoviedb/v3"
	"sync"
)

var (
	client *themoviedb.Client
	lock   = sync.RWMutex{}
)

func Init() error {
	lock.Lock()
	defer lock.Unlock()

	var opts []themoviedb.ClientOptions
	if config.TMDB.Language != "" {
		opts = append(opts, themoviedb.CustomLanguage(config.TMDB.Language))
	}
	if config.TMDB.ApiURL != "" {
		opts = append(opts, themoviedb.CustomAPIURL(config.TMDB.ApiURL))
	}
	if config.TMDB.ImageURL != "" {
		opts = append(opts, themoviedb.CustomImageURL(config.TMDB.ImageURL))
	}
	opts = append(opts, themoviedb.CustomHTTPClient(outbound.GetHTTPClient()))

	var err error
	client, err = themoviedb.NewClient(config.TMDB.ApiKey, opts...)
	return err

}
