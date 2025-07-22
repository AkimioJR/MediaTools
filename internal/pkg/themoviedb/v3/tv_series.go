package themoviedb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TVSeriesDetail struct {
	Vote
	Adult               bool       `json:"adult"`                // 是否成人内容
	BackDropPath        string     `json:"backdrop_path"`        // 背景图片
	CreatedBy           []Creator  `json:"created_by"`           // 创作者列表
	EpisodeRunTime      []uint64   `json:"episode_run_time"`     // 每集时长
	FirstAirDate        string     `json:"first_air_date"`       // 首播日期
	Genres              []Genre    `json:"genres"`               // 类型列表
	Homepage            string     `json:"homepage"`             // 主页
	ID                  uint64     `json:"id"`                   // ID
	InProduction        bool       `json:"in_production"`        // 是否在连载中
	Languages           []string   `json:"languages"`            // 语言列表
	LastAirDate         string     `json:"last_air_date"`        // 最后播出日期
	LastEpisodeToAir    TVEpisode  `json:"last_episode_to_air"`  // 最后一集详情
	Name                string     `json:"name"`                 // 名称
	NextEpisodeToAir    TVEpisode  `json:"next_episode_to_air"`  // 下一集名称
	Networks            []Network  `json:"networks"`             // 播出平台
	NumberOfEpisodes    uint64     `json:"number_of_episodes"`   // 总集
	NumberOfSeasons     uint64     `json:"number_of_seasons"`    // 总季数
	OriginCountry       []string   `json:"origin_country"`       // 原始国家列表
	OriginalLanguage    string     `json:"original_language"`    // 原始语言
	OriginalName        string     `json:"original_name"`        // 原始名称
	Overview            string     `json:"overview"`             // 概述
	Popularity          float64    `json:"popularity"`           // 人气
	PosterPath          string     `json:"poster_path"`          // 海报图片路径
	ProductionCompanies []Company  `json:"production_companies"` // 制作公司列表
	ProductionCountries []Country  `json:"production_countries"` // 制作国家列表
	Seasons             []TVSeason `json:"seasons"`              // 季列表
	SpokenLanguages     []Language `json:"spoken_languages"`     // 语言列表
	Status              string     `json:"status"`               // 状态
	Tagline             string     `json:"tagline"`              // 标语
	Type                string     `json:"type"`                 // 类型
}

// 获取一部电视剧的详细信息。
// Get the details of a TV show.
// https://api.themoviedb.org/3/tv/{series_id}
// https://developer.themoviedb.org/reference/tv-series-details
func (tmdb *TMDB) GetTVSeriesDetails(seriesID uint64, language *string) (*TVSeriesDetail, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	detail := TVSeriesDetail{}
	err := tmdb.DoRequest(http.MethodGet, "/tv/"+strconv.Itoa(int(seriesID)), params, nil, &detail)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧详情失败：%v", err))
	}
	return &detail, nil
}

// 获取已添加到电视剧中的其他标题。
// Get the alternative titles that have been added to a TV show.
// https://api.themoviedb.org/3/tv/{series_id}/alternative_titles
// https://developer.themoviedb.org/reference/tv-series-alternative-titles
//
// country 可选，指定国家(指定一个 ISO-3166-1 值来筛选结果)
func (tmdb *TMDB) GetTVSeriesAlternativeTitles(seriesID uint64, country *string) (*AlternativeTitlesResponse, error) {
	params := url.Values{}
	if country != nil {
		params.Set("country", *country)
	}

	var response AlternativeTitlesResponse
	err := tmdb.DoRequest(http.MethodGet, "/tv/"+strconv.Itoa(int(seriesID))+"/alternative_titles", params, nil, &response)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」别名失败：%v", seriesID, err))
	}
	return &response, nil
}

type TVSeriesGrop struct {
	Description  string    `json:"description"`   // 描述
	EpisodeCount uint64    `json:"episode_count"` // 集数
	ID           uint64    `json:"id"`            // ID
	Name         string    `json:"name"`          // 名称
	Networks     []Network `json:"networks"`      // 播出平台
	Type         uint64    `json:"type"`          // 类型
}

type TVSeriesEpisodeGroupsResponse struct {
	ID      uint64         `json:"id"`      // 电视剧ID
	Results []TVSeriesGrop `json:"results"` // 结果列表
}

// 获取已添加到电视剧中的剧集组。
// Get the episode groups that have been added to a TV show.
// https://api.themoviedb.org/3/tv/{series_id}/episode_groups
// https://developer.themoviedb.org/reference/tv-series-episode-groups
func (tmdb *TMDB) GetTVSeriesEpisodeGroups(seriesID uint64) (*TVSeriesEpisodeGroupsResponse, error) {
	params := url.Values{}
	response := TVSeriesEpisodeGroupsResponse{}
	err := tmdb.DoRequest(http.MethodGet, "/tv/"+strconv.Itoa(int(seriesID))+"/episode_groups", params, nil, &response)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」剧集组失败：%v", seriesID, err))
	}
	return &response, nil
}

// 获取属于某部电视剧的图片。
// Get the images that belong to a TV series.
// https://api.themoviedb.org/3/tv/{series_id}/images
// https://developer.themoviedb.org/reference/tv-series-images
func (tmdb *TMDB) GetTVSeriesImages(seriesID uint64, IncludeImageLanguage *string, language *string) (*Image, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	if IncludeImageLanguage != nil {
		params.Set("include_image_language", *IncludeImageLanguage)
	}

	img := Image{}
	err := tmdb.DoRequest(http.MethodGet, "/tv/"+strconv.Itoa(int(seriesID))+"/images", params, nil, &img)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」图片失败：%v", seriesID, err))
	}
	return &img, nil
}

type TVKeyword struct {
	ID       uint64    `json:"id"`       // 电视剧ID
	Keywords []Keyword `json:"keywords"` // 关键词列表
}

// 获取一部电视剧的关键词列表。
// Get the keywords for a movie by its ID.
// https://api.themoviedb.org/3/tv/{series_id}/keywords
// https://developer.themoviedb.org/reference/movie-keywords
func (tmdb *TMDB) TVSeriesKeywords(seriesID uint64) (*TVKeyword, error) {
	kyword := TVKeyword{}
	err := tmdb.DoRequest(http.MethodGet, "/tv/"+strconv.Itoa(int(seriesID))+"/keywords", url.Values{}, nil, &kyword)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」关键词失败：%v", seriesID, err))
	}
	if kyword.ID == 0 {
		return nil, NewTMDBError(nil, fmt.Sprintf("电视剧「%d」不存在或没有关键词", seriesID))
	}
	return &kyword, nil
}
