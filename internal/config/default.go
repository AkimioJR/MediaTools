package config

import "MediaTools/internal/schemas/storage"

var defaultConfig = Configuration{
	Log: LogConfig{
		ConsoleLevel: "info",
		FileLevel:    "debug",
		FileDir:      "logs",
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
	Storages: []StorageConfig{
		{
			Type: storage.StorageLocal, // 默认使用本地存储
			Data: map[string]string{},
		},
	},
	Media: MediaConfig{
		Format: FormatConfig{
			Movie: "{{.Title}} ({{.Year}})/{{.Title}} ({{.Year}}){{if .Part}} -{{.Part}}{{end}}{{if .Version}} -v{{.Version}}{{end}}{{if .ReleaseGroups}} -{{end}}{{range .ReleaseGroups}}@{{.}}{{end}}{{if .ResourcePix}} -{{.ResourcePix}}{{end}}{{if .ResourceType}} -{{.ResourceType}}{{end}}{{if .ResourceEffect}} -{{end}}{{range .ResourceEffect}}@{{.}}{{end}}{{if .Platform}} -{{.Platform}}{{end}}{{if .VideoEncode}} -{{.VideoEncode}}{{end}}{{if .AudioEncode}} -{{.AudioEncode}}{{end}}{{.FileExtension}}",
			TV:    "{{.Title}} ({{.Year}})/Season {{.Season}}/{{.Title}} {{.SeasonStr}}{{.EpisodeStr}}{{if .EpisodeTitle}} {{.EpisodeTitle}}{{end}}{{if .Part}} -{{.Part}}{{end}}{{if .Version}} -v{{.Version}}{{end}}{{if .ReleaseGroups}} -{{end}}{{range .ReleaseGroups}}@{{.}}{{end}}{{if .ResourcePix}} -{{.ResourcePix}}{{end}}{{if .ResourceType}} -{{.ResourceType}}{{end}}{{if .ResourceEffect}} -{{end}}{{range .ResourceEffect}}@{{.}}{{end}}{{if .Platform}} -{{.Platform}}{{end}}{{if .VideoEncode}} -{{.VideoEncode}}{{end}}{{if .AudioEncode}} -{{.AudioEncode}}{{end}}{{.FileExtension}}",
		},
	},
}
