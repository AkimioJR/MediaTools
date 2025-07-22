package themoviedb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type MoviesDetails struct {
	Vote
	Adult               bool       `json:"adult"`                 // 是否为成人内容
	BackDropPath        string     `json:"backdrop_path"`         // 背景图片路径
	BelongsToCollection Collection `json:"belongs_to_collection"` // 所属系列
	Budget              uint64     `json:"budget"`                // 预算
	Genres              []Genre    `json:"genres"`                // 类型ID列表
	Homepage            string     `json:"homepage"`              // 主页
	ID                  uint64     `json:"id"`                    // ID
	OriginalLanguage    string     `json:"original_language"`     // 原始语言
	OriginalTitle       string     `json:"original_title"`        // 原始标题
	Overview            string     `json:"overview"`              // 概述
	Popularity          float64    `json:"popularity"`            // 人气
	PosterPath          string     `json:"poster_path"`           // 海报图片路径
	ProductionCompanies []Company  `json:"production_companies"`  // 制作公司列表
	ReleaseDate         string     `json:"release_date"`          // 发行日期
	Revenue             uint64     `json:"revenue"`               // 收入
	Runtime             uint64     `json:"runtime"`               // 运行时长
	SpokenLanguages     []Language `json:"spoken_languages"`      // 语言列表
	Status              string     `json:"status"`                // 状态
	Tagline             string     `json:"tagline"`               // 标语
	Title               string     `json:"title"`                 // 标题
	Video               bool       `json:"video"`                 // 是否是视频
}

// 通过ID获取一部电影的顶级详情。
// Get the top level details of a movie by ID.
// https://api.themoviedb.org/3/movie/{movie_id}
// https://developer.themoviedb.org/reference/movie-details
func (tmdb *TMDB) GetMovieDetails(movieID uint64, language *string) (*MoviesDetails, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	resp := MoviesDetails{}
	err := tmdb.DoRequest(http.MethodGet, "/movie/"+strconv.Itoa(int(movieID)), params, nil, &resp)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影详情失败：%v", err))
	}
	return &resp, nil
}

// 获取属于某部电影的图片。
// Get the images that belong to a movie.
// https://api.themoviedb.org/3/movie/{movie_id}/images
// https://developer.themoviedb.org/reference/movie-images
func (tmdb *TMDB) GetMovieImages(movieID uint64, language *string, IncludeImageLanguage *string) (*Image, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	if IncludeImageLanguage != nil {
		params.Set("include_image_language", *IncludeImageLanguage)
	}

	resp := Image{}
	err := tmdb.DoRequest(http.MethodGet, "/movie/"+strconv.Itoa(int(movieID))+"/images", params, nil, &resp)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影「%d」图片失败：%v", movieID, err))
	}
	return &resp, nil
}

type MovieKeywordsResponse struct {
	ID       uint64    `json:"id"`       // 电影ID
	Keywords []Keyword `json:"keywords"` // 关键词列表
}

// 获取一部电影的关键词列表。
// Get the keywords for a movie by its ID.
// https://api.themoviedb.org/3/movie/{movie_id}/keywords
// https://developer.themoviedb.org/reference/movie-keywords
func (tmdb *TMDB) GetMovieKeywords(movieID uint64) (*MovieKeywordsResponse, error) {
	resp := MovieKeywordsResponse{}
	err := tmdb.DoRequest(http.MethodGet, "/movie/"+strconv.Itoa(int(movieID))+"/keywords", url.Values{}, nil, &resp)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影「%d」关键词失败：%v", movieID, err))
	}
	if resp.ID == 0 {
		return nil, NewTMDBError(nil, fmt.Sprintf("电影「%d」不存在或没有关键词", movieID))
	}
	return &resp, nil
}
