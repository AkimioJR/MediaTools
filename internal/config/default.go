package config

import "MediaTools/internal/pkg/loghook"

var defaultConfig = Configuration{
	Log: LogConfig{
		Level: loghook.LevelInfo,
		Path:  "logs",
	},
	TMDB: TMDBConfig{
		ApiKey:   "YOUR_TMDB_API_KEY", // 请替换为您的 TMDB API Key
		ApiURL:   "https://api.themoviedb.org",
		ImageURL: "https://image.tmdb.org",
	},
	Fanart: FanartConfig{
		ApiKey: "YOUR_FANART_API_KEY", // 请替换为您的 Fanart API Key
		ApiURL: "https://webservice.fanart.tv",
	},
}
