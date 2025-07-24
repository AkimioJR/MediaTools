package schemas

import (
	"MediaTools/encode"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/pkg/themoviedb/v3"
	"fmt"
	"html/template"
	"path"
	"strconv"
	"strings"
)

type TMDBTVInfo struct {
	SeriesInfo    *themoviedb.TVSerieDetail   // TMDB 电视剧相关信息
	SeasonInfo    *themoviedb.TVSeasonDetail  // TMDB 电视剧季相关信息
	EpisodeInfo   *themoviedb.TVEpisodeDetail // TMDB 电视剧集相关信息
	SeasonNumber  int                         // 季数
	EpisodeNumber int                         // 集数
}

type TMDBInfo struct {
	MovieInfo *themoviedb.MovieDetail // TMDB 电影相关信息
	TVInfo    TMDBTVInfo              // TMDB 电视剧相关信息
}

// 用于记录数据库中相关数据信息
type MediaInfo struct {
	MediaType meta.MediaType // 电影、电视剧等
	TMDBID    int            // TMDB ID
	TMDBInfo  TMDBInfo       // TMDB 相关信息

	// IMDBID      string // IMDb ID
	// TVDBID      uint64 // TVDB ID
	// DoubanID    string // 豆瓣 ID
	// BangumiID   string // 番组计划 ID
	// DoubanInfo  any    // 豆瓣相关信息
	// BangumiInfo any    // 番组计划相关信息
}

// 用于媒体库整理重命名可选字段模板
type MediaItem struct {
	Title         string                 `json:"title"`          // 标题
	OriginalTitle string                 `json:"original_title"` // 原始标题
	Year          int                    `json:"year"`           // 年份
	MediaType     meta.MediaType         `json:"media_type"`     // 电影、电视剧
	TMDBID        int                    `json:"tmdb_id"`        // TMDB ID
	Part          string                 `json:"part"`           // 分段
	Version       uint8                  `json:"version"`        // 版本号
	ReleaseGroups []string               `json:"release_groups"` // 发布组
	Platform      meta.StreamingPlatform `json:"platform"`       // 流媒体平台
	FileExtension string                 `json:"file_extension"` // 文件扩展名

	// 资源相关信息
	ResourceType   meta.ResourceType     `json:"resource_type"`   // 资源类型
	ResourceEffect []meta.ResourceEffect `json:"resource_effect"` // 资源效果
	ResourcePix    meta.ResourcePix      `json:"resource_pix"`    // 分辨率
	VideoEncode    encode.VideoEncode    `json:"video_encode"`    // 视频编码
	AudioEncode    encode.AudioEncode    `json:"audio_encode"`    // 音频编码

	// 电视剧数据
	Season       int    `json:"season"`        // 季数 -1表示无季数
	SeasonStr    string `json:"season_str"`    // 季 S01 S01-S03
	SeasonYear   int    `json:"season_year"`   // 季年份
	Episode      int    `json:"episode"`       // 集数 -1表示无集数
	EpisodeStr   string `json:"episode_str"`   // 集 E12 E12-E15
	EpisodeTitle string `json:"episode_title"` // 集标题
	EpisodeDate  string `json:"episode_date"`  // 集发布日期
}

func NewMediaItem(videoMeta *meta.VideoMeta, info *MediaInfo) (*MediaItem, error) {
	item := MediaItem{
		MediaType:      info.MediaType,
		TMDBID:         info.TMDBID,
		Part:           videoMeta.Part,
		Version:        videoMeta.Version,
		ReleaseGroups:  videoMeta.ReleaseGroups,
		Platform:       videoMeta.Platform,
		ResourceType:   videoMeta.ResourceType,
		ResourceEffect: videoMeta.ResourceEffect,
		ResourcePix:    videoMeta.ResourcePix,
		VideoEncode:    videoMeta.VideoEncode,
		AudioEncode:    videoMeta.AudioEncode,
		FileExtension:  path.Ext(videoMeta.OrginalTitle),
	}

	switch info.MediaType {
	case meta.MediaTypeMovie:
		item.Title = info.TMDBInfo.MovieInfo.Title
		item.OriginalTitle = info.TMDBInfo.MovieInfo.OriginalTitle
		year, err := strconv.Atoi(info.TMDBInfo.MovieInfo.ReleaseDate[:4])
		if err == nil {
			item.Year = year
		}
	case meta.MediaTypeTV:
		item.Title = info.TMDBInfo.TVInfo.SeriesInfo.Name
		item.OriginalTitle = info.TMDBInfo.TVInfo.SeriesInfo.OriginalName
		year, err := strconv.Atoi(info.TMDBInfo.TVInfo.SeriesInfo.FirstAirDate[:4])
		if err == nil {
			item.Year = year
		}
		item.Season = videoMeta.Season
		item.SeasonStr = videoMeta.GetSeasonStr()
		year, err = strconv.Atoi(info.TMDBInfo.TVInfo.SeasonInfo.AirDate[:4])
		if err == nil {
			item.SeasonYear = year
		}
		item.Episode = videoMeta.Episode
		item.EpisodeStr = videoMeta.GetEpisodeStr()
		item.EpisodeTitle = info.TMDBInfo.TVInfo.EpisodeInfo.Name
		item.EpisodeDate = info.TMDBInfo.TVInfo.EpisodeInfo.AirDate
	default:
		return nil, fmt.Errorf("不支持的媒体类型: %s", info.MediaType.String())
	}

	return &item, nil
}

func (item *MediaItem) Format(format string) (string, error) {
	tmpl, err := template.New("mediaName").Parse(format)
	if err != nil {
		return "", fmt.Errorf("解析模板字符串「%s」失败: %v", format, err)
	}
	var buffer strings.Builder
	if err := tmpl.Execute(&buffer, item); err != nil {
		return "", fmt.Errorf("渲染模板失败: %v", err)
	}
	return buffer.String(), nil
}
