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

type FanartConfig struct {
	APIKey   string   `yaml:"api_key"`
	Language []string `yaml:"language"` // 语言顺序
}

type Configuration struct {
	Log          LogConfig          `yaml:"log"`
	TMDB         TMDBConfig         `yaml:"tmdb"`
	Fanart       FanartConfig       `yaml:"fanart"`
	MediaLibrary MediaLibraryConfig `yaml:"media_library"`
}

type MediaLibraryConfig struct {
	Libraries   []string `yaml:"libraries"`    // 媒体库路径列表
	MovieFormat string   `yaml:"movie_format"` // 电影格式
	TVFormat    string   `yaml:"tv_format"`    // 电视剧格式
}
