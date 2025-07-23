package config

type LogConfig struct {
	Level string `yaml:"level"` // 日志级别
	Path  string `yaml:"path"`  // 日志文件目录
}

type TMDBConfig struct {
	APIURL   string `yaml:"api_url"`
	ImageURL string `yaml:"image_url"`
	APIKey   string `yaml:"api_key"`
	Language string `yaml:"language"`
}

type Configuration struct {
	Log  LogConfig  `yaml:"log"`
	TMDB TMDBConfig `yaml:"tmdb"`
}
