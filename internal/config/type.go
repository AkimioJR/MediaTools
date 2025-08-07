package config

import "MediaTools/internal/schemas"

type LogConfig struct {
	Level string `json:"level"` // 日志级别
	Path  string `json:"path"`  // 日志文件目录
}

type TMDBConfig struct {
	ApiURL               string `json:"api_url"`                // TMDB API URL
	ImageURL             string `json:"image_url"`              // 图片 API URL
	ApiKey               string `json:"api_key"`                // API Key
	Language             string `json:"language"`               // 语言
	IncludeImageLanguage string `json:"include_image_language"` // 包含的图片语言
}

type FanartConfig struct {
	ApiKey    string   `json:"api_key"`
	ApiURL    string   `json:"api_url"`
	Languages []string `json:"languages"` // 语言顺序
}

type Configuration struct {
	Log    LogConfig    `json:"log"`
	TMDB   TMDBConfig   `json:"tmdb"`
	Fanart FanartConfig `json:"fanart"`
	Media  MediaConfig  `json:"media"`
}

type MediaConfig struct {
	Libraries  []LibraryConfig  `json:"libraries"`   // 媒体库路径列表
	Format     FormatConfig     `json:"format"`      // 媒体格式配置
	CustomWord CustomWordConfig `json:"custom_word"` // 自定义识别词配置
}

type FormatConfig struct {
	Movie string `json:"movie"` // 电影格式
	TV    string `json:"tv"`    // 电视剧格式
}

type LibraryConfig struct {
	Name               string               `json:"name"`                 // 媒体库名称
	SrcPath            string               `json:"src_path"`             // 源路径
	SrcType            schemas.StorageType  `json:"src_type"`             // 源类型
	DstType            schemas.StorageType  `json:"dst_type"`             // 目标类型
	DstPath            string               `json:"dst_path"`             // 目标路径
	TransferType       schemas.TransferType `json:"transfer_type"`        // 传输类型
	OrganizeByType     bool                 `json:"organize_by_type"`     // 是否按类型分文件夹
	OrganizeByCategory bool                 `json:"organize_by_category"` // 是否按分类分文件夹
}

type CustomWordConfig struct {
	IdentifyWord  []string `json:"identify_word"` // 自定义识别词
	Customization []string `json:"customization"` // 自定义占位置词
	ExcludeWords  []string `json:"exclude_words"` // 自定义排除词
}
