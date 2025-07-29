package tmdb_controller

import (
	"MediaTools/internal/config"
	"MediaTools/internal/pkg/themoviedb/v3"
)

var client *themoviedb.TMDB

func Init() error {
	var otps []themoviedb.TMDBOptions
	if config.TMDB.Language != "" {
		otps = append(otps, themoviedb.CustomLanguage(config.TMDB.Language))
	}
	if config.TMDB.ApiURL != "" {
		otps = append(otps, themoviedb.CustomAPIURL(config.TMDB.ApiURL))
	}
	if config.TMDB.ImageURL != "" {
		otps = append(otps, themoviedb.CustomImageURL(config.TMDB.ImageURL))
	}
	client = themoviedb.NewTMDB(config.TMDB.ApiKey, otps...)
	return nil

}
