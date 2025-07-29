package tmdb_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/themoviedb/v3"
	"sync"
)

var (
	client *themoviedb.TMDB
	lock   = sync.RWMutex{}
)

func Init() error {
	lock.Lock()
	defer lock.Unlock()

	var opts []themoviedb.TMDBOptions
	if config.TMDB.Language != "" {
		opts = append(opts, themoviedb.CustomLanguage(config.TMDB.Language))
	}
	if config.TMDB.ApiURL != "" {
		opts = append(opts, themoviedb.CustomAPIURL(config.TMDB.ApiURL))
	}
	if config.TMDB.ImageURL != "" {
		opts = append(opts, themoviedb.CustomImageURL(config.TMDB.ImageURL))
	}

	client = themoviedb.NewTMDB(config.TMDB.ApiKey, opts...)
	return nil

}
