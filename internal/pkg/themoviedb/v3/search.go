package themoviedb

import (
	"MediaTools/utils"
	"fmt"
	"net/http"
)

type SearchResponse[T any] struct {
	Page         int `json:"page"`          // 当前页码
	TotalPages   int `json:"total_pages"`   // 总页数
	Result       []T `json:"results"`       // 搜索结果
	TotalResults int `json:"total_results"` // 总结果数
}

// ========= 请求字段 =========

type SearchKeywordParams struct {
	Query string  `query:"query"` // 搜索关键词
	Page  *uint32 `query:"page"`  // 页码
}

type SearchCompanyParams SearchKeywordParams

type SearchMultiParams struct {
	SearchKeywordParams
	IncludeAdult *bool   `query:"include_adult"` // 是否包含成人内容
	Language     *string `query:"language"`      // 语言
}

type SearchCollectionParams struct {
	SearchMultiParams
	Region *string `query:"region"` // 地区
}

type SearchMovieParams struct {
	SearchCollectionParams
	PrimaryReleaseYear *string `query:"primary_release_year"` // 主发行年份
	Year               *string `query:"year"`                 // 年份
}

type SearchTVSParams struct {
	SearchMultiParams
	FirstAirDateYear *uint32 `query:"first_air_date_year"` // 首播年份
	Year             *uint32 `query:"year"`                // 年份
}

type SearchPersonParams SearchMultiParams

// ========= 响应字段 =========

type SearchKeywordResponse Keyword

type SearchCompanyResponse struct {
	SearchKeywordResponse
	LogoPath      string `json:"logo_path"`      // 公司Logo路径
	OriginCountry string `json:"origin_country"` // 原始国家
}

type SearchCollectionResponse struct {
	SearchKeywordResponse
	Adult            bool   `json:"adult"`             // 是否成人内容
	BackDropPath     string `json:"backdrop_path"`     // 背景图片路径
	OriginalLanguage string `json:"original_language"` // 原始语言
	OriginalName     string `json:"original_name"`     // 原始名称
	Overview         string `json:"overview"`          // 概述
	PosterPath       string `json:"poster_path"`       // 海报图片路径
}

type BaseMediaResponse struct {
	Vote
	Adult            bool    `json:"adult"`             // 是否成人内容
	BackDropPath     string  `json:"backdrop_path"`     // 背景图片路径
	GenreIDs         []int   `json:"genre_ids"`         // 类型ID列表
	ID               int     `json:"id"`                // ID
	OriginalLanguage string  `json:"original_language"` // 原始语言
	Overview         string  `json:"overview"`          // 概述
	Popularity       float64 `json:"popularity"`        // 人气
	PosterPath       string  `json:"poster_path"`       // 海报图片路径
}

type SearchMovieResponse struct {
	BaseMediaResponse
	OriginalTitle string `json:"original_title"` // 原始标题
	Title         string `json:"title"`          // 标题
	ReleaseDate   string `json:"release_date"`   // 发行日期
	Video         bool   `json:"video"`          // 是否视频
}

type SearchTVResponse struct {
	BaseMediaResponse
	OriginalCountry []string `json:"original_country"` // 原始国家
	OriginalName    string   `json:"original_name"`    // 原始名称
	Name            string   `json:"name"`             // 名称
	FirstAirDate    string   `json:"first_air_date"`   // 首播日期
}

type SearchMultiResponse struct {
	BaseMediaResponse
	MediaType     string `json:"media_type"`     // 媒体类型
	OriginalTitle string `json:"original_title"` // 原始标题（电影）
	Title         string `json:"title"`          // 标题（电影）
	ReleaseDate   string `json:"release_date"`   // 发行日期（电影）
	Video         bool   `json:"video"`          // 是否视频（电影）
}

type SearchPersonResponse struct {
	Adult              bool                  `json:"adult"`                // 是否成人内容
	Gender             int                   `json:"gender"`               // 性别
	ID                 int                   `json:"id"`                   // ID
	KnownForDepartment string                `json:"known_for_department"` // 知名部门
	Name               string                `json:"name"`                 // 姓名
	OriginalName       string                `json:"original_name"`        // 原始姓名
	Popularity         float64               `json:"popularity"`           // 人气
	ProfilePath        string                `json:"profile_path"`         // 个人资料图片路径
	KnownFor           []SearchMultiResponse `json:"known_for"`            // 知名作品
}

// 按收藏的原名、译名及别名进行搜索。
// Search for collections by their original, translated and alternative names.
// https://developer.themoviedb.org/reference/search-collection
func (tmdb *TMDB) SearchCollection(params SearchCollectionParams) (*SearchResponse[SearchCollectionResponse], error) {
	var resp SearchResponse[SearchCollectionResponse]
	if params.Language == nil {
		params.Language = &tmdb.language
	}
	err := tmdb.DoRequest(
		http.MethodGet,
		"/search/collection",
		utils.StructToQuery(params),
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("搜索合集「%s」失败：%v", params.Query, err))
	}
	return &resp, nil
}

// 按公司的原名和别名搜索。
// Search for companies by their original and alternative names.
// https://developer.themoviedb.org/reference/search-company
func (tmdb *TMDB) SearchCompany(params SearchCompanyParams) ([]SearchCompanyResponse, error) {
	var resp SearchResponse[SearchCompanyResponse]
	err := tmdb.DoRequest(
		http.MethodGet,
		"/search/company",
		utils.StructToQuery(params),
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("搜索公司「%s」失败：%v", params.Query, err))
	}
	return resp.Result, nil
}

// 按名称搜索关键词。
// Search for keywords by their name.
// https://developer.themoviedb.org/reference/search-keyword
func (tmdb *TMDB) SearchKeyword(params SearchKeywordParams) (*SearchResponse[SearchKeywordResponse], error) {
	var resp SearchResponse[SearchKeywordResponse]
	err := tmdb.DoRequest(
		http.MethodGet,
		"/search/keyword",
		utils.StructToQuery(params),
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("搜索关键词「%s」失败：%v", params.Query, err))
	}
	return &resp, nil
}

// 按电影的原名、译名和别名搜索。
// Search for movies by their original, translated and alternative titles.
// https://developer.themoviedb.org/reference/search-movie
func (tmdb *TMDB) SearchMovie(params SearchMovieParams) (*SearchResponse[SearchMovieResponse], error) {
	var resp SearchResponse[SearchMovieResponse]
	if params.Language == nil {
		params.Language = &tmdb.language
	}
	err := tmdb.DoRequest(
		http.MethodGet,
		"/search/movie", utils.StructToQuery(params),
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("搜索电影「%s」失败：%v", params.Query, err))
	}
	return &resp, nil
}

// 当你想在一次请求中搜索电影、电视节目和人物时，请使用多重搜索。
// Use multi search when you want to search for movies, TV shows and people in a single request.
// https://developer.themoviedb.org/reference/search-multi
func (tmdb *TMDB) SearchMulti(params SearchMultiParams) (*SearchResponse[SearchMultiResponse], error) {
	var resp SearchResponse[SearchMultiResponse]
	if params.Language == nil {
		params.Language = &tmdb.language
	}
	err := tmdb.DoRequest(
		http.MethodGet,
		"/search/multi",
		utils.StructToQuery(params),
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("搜索多种类型「%s」失败：%v", params.Query, err))
	}
	return &resp, nil
}

// 按人名及其曾用名搜索人物。
// Search for people by their name and also known as names.
// https://developer.themoviedb.org/reference/search-person
func (tmdb *TMDB) SearchPerson(params SearchPersonParams) (*SearchResponse[SearchPersonResponse], error) {
	var resp SearchResponse[SearchPersonResponse]
	if params.Language == nil {
		params.Language = &tmdb.language
	}
	err := tmdb.DoRequest(
		http.MethodGet,
		"/search/person",
		utils.StructToQuery(params),
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("搜索人物「%s」失败：%v", params.Query, err))
	}
	return &resp, nil
}

// 按电视节目的原名、译名及别名搜索。
// Search for TV shows by their original, translated and also known as names.
// https://developer.themoviedb.org/reference/search-tv
func (tmdb *TMDB) SearchTV(params SearchTVSParams) (*SearchResponse[SearchTVResponse], error) {
	var resp SearchResponse[SearchTVResponse]
	if params.Language == nil {
		params.Language = &tmdb.language
	}
	err := tmdb.DoRequest(
		http.MethodGet,
		"/search/tv",
		utils.StructToQuery(params),
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("搜索电视剧「%s」失败：%v", params.Query, err))
	}
	return &resp, nil
}
