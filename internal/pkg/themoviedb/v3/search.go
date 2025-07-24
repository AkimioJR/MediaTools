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

type SearchCompanyParams struct {
	Query string  `query:"query"` // 搜索关键词
	Page  *uint32 `query:"page"`  // 页码
}

type SearchMultiParams struct {
	Query        string  `query:"query"`         // 搜索关键词
	Page         *uint32 `query:"page"`          // 页码
	IncludeAdult *bool   `query:"include_adult"` // 是否包含成人内容
	Language     *string `query:"language"`      // 语言
}

type SearchCollectionParams struct {
	Query        string  `query:"query"`         // 搜索关键词
	Page         *uint32 `query:"page"`          // 页码
	IncludeAdult *bool   `query:"include_adult"` // 是否包含成人内容
	Language     *string `query:"language"`      // 语言
	Region       *string `query:"region"`        // 地区
}

type SearchMovieParams struct {
	Query              string  `query:"query"`                // 搜索关键词
	Page               *uint32 `query:"page"`                 // 页码
	IncludeAdult       *bool   `query:"include_adult"`        // 是否包含成人内容
	Language           *string `query:"language"`             // 语言
	Region             *string `query:"region"`               // 地区
	PrimaryReleaseYear *string `query:"primary_release_year"` // 主发行年份
	Year               *string `query:"year"`                 // 年份
}

type SearchTVSParams struct {
	Query            string  `query:"query"`               // 搜索关键词
	Page             *uint32 `query:"page"`                // 页码
	IncludeAdult     *bool   `query:"include_adult"`       // 是否包含成人内容
	Language         *string `query:"language"`            // 语言
	FirstAirDateYear *uint32 `query:"first_air_date_year"` // 首播年份
	Year             *uint32 `query:"year"`                // 年份
}

type SearchPersonParams struct {
	Query        string  `query:"query"`         // 搜索关键词
	Page         *uint32 `query:"page"`          // 页码
	IncludeAdult *bool   `query:"include_adult"` // 是否包含成人内容
	Language     *string `query:"language"`      // 语言
}

// ========= 响应字段 =========

type SearchKeywordResponse struct {
	ID   int    `json:"id"`   // 关键词ID
	Name string `json:"name"` // 关键词名称
}

type SearchCompanyResponse struct {
	ID            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type SearchCollectionResponse struct {
	Adult            bool   `json:"adult"`
	BackdropPath     string `json:"backdrop_path"`
	ID               int    `json:"id"`
	Name             string `json:"name"`
	OriginalLanguage string `json:"original_language"`
	OriginalName     string `json:"original_name"`
	Overview         string `json:"overview"`
	PosterPath       string `json:"poster_path"`
}

type SearchMovieResponse struct {
	Adult            bool    `json:"adult"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIds         []int   `json:"genre_ids"`
	ID               int     `json:"id"`
	OriginalLanguage string  `json:"original_language"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	Popularity       float64 `json:"popularity"`
	PosterPath       string  `json:"poster_path"`
	ReleaseDate      string  `json:"release_date"`
	Title            string  `json:"title"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type SearchTVResponse struct {
	Adult            bool     `json:"adult"`
	BackdropPath     string   `json:"backdrop_path"`
	GenreIds         []int    `json:"genre_ids"`
	ID               int      `json:"id"`
	OriginCountry    []string `json:"origin_country"`
	OriginalLanguage string   `json:"original_language"`
	OriginalName     string   `json:"original_name"`
	Overview         string   `json:"overview"`
	Popularity       float64  `json:"popularity"`
	PosterPath       string   `json:"poster_path"`
	FirstAirDate     string   `json:"first_air_date"`
	Name             string   `json:"name"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
}

type SearchMultiResponse struct {
	Adult            bool     `json:"adult"`
	BackdropPath     string   `json:"backdrop_path"`
	ID               int      `json:"id"`
	Title            string   `json:"title,omitempty"`
	OriginalLanguage string   `json:"original_language"`
	OriginalTitle    string   `json:"original_title,omitempty"`
	Overview         string   `json:"overview"`
	PosterPath       string   `json:"poster_path"`
	MediaType        string   `json:"media_type"`
	GenreIds         []int    `json:"genre_ids"`
	Popularity       float64  `json:"popularity"`
	ReleaseDate      string   `json:"release_date,omitempty"`
	Video            bool     `json:"video,omitempty"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
	Name             string   `json:"name,omitempty"`
	OriginalName     string   `json:"original_name,omitempty"`
	FirstAirDate     string   `json:"first_air_date,omitempty"`
	OriginCountry    []string `json:"origin_country,omitempty"`
}

type SearchPersonResponse struct {
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	ID                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float64 `json:"popularity"`
	ProfilePath        string  `json:"profile_path"`
	KnownFor           []struct {
		Adult            bool    `json:"adult"`
		BackdropPath     string  `json:"backdrop_path"`
		ID               int     `json:"id"`
		Title            string  `json:"title"`
		OriginalLanguage string  `json:"original_language"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		PosterPath       string  `json:"poster_path"`
		MediaType        string  `json:"media_type"`
		GenreIds         []int   `json:"genre_ids"`
		Popularity       float64 `json:"popularity"`
		ReleaseDate      string  `json:"release_date"`
		Video            bool    `json:"video"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
	} `json:"known_for"`
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
