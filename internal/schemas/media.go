package schemas

import (
	"MediaTools/encode"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/pkg/themoviedb/v3"
	"fmt"
	"path/filepath"
	"strconv"
)

type TMDBTVInfo struct {
	SerieInfo     *themoviedb.TVSerieDetail   // TMDB 电视剧相关信息
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

	IMDBID string // IMDb ID
	TVDBID int    // TVDB ID

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
	Part          string                 `json:"part"`           // 分段
	Version       uint8                  `json:"version"`        // 版本号
	ReleaseGroups []string               `json:"release_groups"` // 发布组
	Platform      meta.StreamingPlatform `json:"platform"`       // 流媒体平台
	FileExtension string                 `json:"file_extension"` // 文件扩展名

	Customization []string `json:"customization"` // 自定义词

	// ID 信息
	TMDBID int    `json:"tmdb_id"` // TMDB ID
	IMDBID string `json:"imdb_id"` // IMDb ID
	TVDBID int    `json:"tvdb_id"` // TVDB ID

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
		Part:           videoMeta.Part,
		Version:        videoMeta.Version,
		ReleaseGroups:  videoMeta.ReleaseGroups,
		Platform:       videoMeta.Platform,
		ResourceType:   videoMeta.ResourceType,
		ResourceEffect: videoMeta.ResourceEffect,
		ResourcePix:    videoMeta.ResourcePix,
		VideoEncode:    videoMeta.VideoEncode,
		AudioEncode:    videoMeta.AudioEncode,
		FileExtension:  filepath.Ext(videoMeta.OrginalTitle),
		Customization:  videoMeta.Customization,

		TMDBID: info.TMDBID,
		IMDBID: info.IMDBID,
		TVDBID: info.TVDBID,
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
		item.Title = info.TMDBInfo.TVInfo.SerieInfo.Name
		item.OriginalTitle = info.TMDBInfo.TVInfo.SerieInfo.OriginalName
		year, err := strconv.Atoi(info.TMDBInfo.TVInfo.SerieInfo.FirstAirDate[:4])
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
