package tmdb_controller

import (
	"MediaTools/internal/pkg/themoviedb/v3"
)

var client *themoviedb.TMDB

func Init(apikey string, opts ...themoviedb.TMDBOptions) {
	client = themoviedb.NewTMDB(apikey, opts...)
}
