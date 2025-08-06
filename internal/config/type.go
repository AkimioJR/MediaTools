package config

import "MediaTools/internal/schemas"

type LogConfig struct {
	Level string `yaml:"level"` // 日志级别
	Path  string `yaml:"path"`  // 日志文件目录
}

type TMDBConfig struct {
	ApiURL   string `yaml:"api_url"`
	ImageURL string `yaml:"image_url"`
	ApiKey   string `yaml:"api_key"`
	Language string `yaml:"language"`
}

type FanartConfig struct {
	ApiKey    string   `yaml:"api_key"`
	ApiURL    string   `yaml:"api_url"`
	Languages []string `yaml:"languages"` // 语言顺序
}

type Configuration struct {
	Log    LogConfig    `yaml:"log"`
	TMDB   TMDBConfig   `yaml:"tmdb"`
	Fanart FanartConfig `yaml:"fanart"`
	Media  MediaConfig  `yaml:"media"`
}

type MediaConfig struct {
	Libraries []LibraryConfig `yaml:"libraries"` // 媒体库路径列表
	Format    FormatConfig    `yaml:"format"`    // 媒体格式配置
}

type FormatConfig struct {
	Movie string `yaml:"movie"` // 电影格式
	TV    string `yaml:"tv"`    // 电视剧格式
}

type LibraryConfig struct {
	Name               string               `yaml:"name"`                 // 媒体库名称
	SrcPath            string               `yaml:"src_path"`             // 源路径
	SrcType            schemas.StorageType  `yaml:"src_type"`             // 源类型
	DstType            schemas.StorageType  `yaml:"dst_type"`             // 目标类型
	DstPath            string               `yaml:"dst_path"`             // 目标路径
	TransferType       schemas.TransferType `yaml:"transfer_type"`        // 传输类型
	OrganizeByType     bool                 `yaml:"organize_by_type"`     // 是否按类型分文件夹
	OrganizeByCategory bool                 `yaml:"organize_by_category"` // 是否按分类分文件夹
}
