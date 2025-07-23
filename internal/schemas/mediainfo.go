package schemas

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/pkg/themoviedb/v3"
)

type TMDBTVInfo struct {
	SeriesInfo  *themoviedb.TVSeriesDetail  // TMDB 电视剧相关信息
	SeasonInfo  *themoviedb.TVSeasonDetail  // TMDB 电视剧季相关信息
	EpisodeInfo *themoviedb.TVEpisodeDetail // TMDB 电视剧集相关信息
}

type TMDBInfo struct {
	MovieInfo *themoviedb.MovieDetail // TMDB 电影相关信息
	TVInfo    TMDBTVInfo              // TMDB 电视剧相关信息
}

type MediaInfo struct {
	MediaType meta.MediaType // 电影、电视剧等
	TMDBID    uint64         // TMDB ID
	TMDBInfo  TMDBInfo       // TMDB 相关信息

	// IMDBID      string // IMDb ID
	// TVDBID      uint64 // TVDB ID
	// DoubanID    string // 豆瓣 ID
	// BangumiID   string // 番组计划 ID
	// DoubanInfo  any    // 豆瓣相关信息
	// BangumiInfo any    // 番组计划相关信息
}
