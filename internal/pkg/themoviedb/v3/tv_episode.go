package themoviedb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TVEpisodeDetail struct {
	AirDate        string      `json:"air_date"`        // 首播日期
	EpisodeNumber  int         `json:"episode_number"`  // 集数
	EpisodeType    string      `json:"episode_type"`    // 集类型
	ID             int         `json:"id"`              // 集 ID
	Name           string      `json:"name"`            // 名称
	Overview       string      `json:"overview"`        // 概述
	ProductionCode string      `json:"production_code"` // 制作代码
	Runtime        int         `json:"runtime"`         // 时长
	SeasonNumber   int         `json:"season_number"`   // 季数
	ShowID         int         `json:"show_id"`         // 电视剧 ID
	StillPath      string      `json:"still_path"`      // 静态图片路径
	VoteAverage    float64     `json:"vote_average"`    // 平均评分
	VoteCount      int         `json:"vote_count"`      // 投票数
	Crew           []Crew      `json:"crew"`            // 工作人员列表
	GuestStars     []GuestStar `json:"guest_stars"`     // 特邀演员列表
}

// 查询电视剧单集的详细信息。
// Query the details of a TV episode.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/episode/{episode_number}
// https://developer.themoviedb.org/reference/tv-episode-details
func (tmdb *TMDB) GetTVEpisodeDetail(seriesID int, seasonNumber int, episodeNumber int, language *string) (*TVEpisodeDetail, error) {
	var resp TVEpisodeDetail

	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(seriesID))+"/season/"+strconv.Itoa(int(seasonNumber))+"/episode/"+strconv.Itoa(int(episodeNumber)),
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季 第 %d 集」详情失败：%v", seriesID, seasonNumber, episodeNumber, err))
	}
	return &resp, nil
}

// 获取一部电视剧单集的演员列表和工作人员列表。
// Get the cast and crew for a TV episode by its ID.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/episode/{episode_number}/credits
// https://developer.themoviedb.org/reference/tv-episode-credits
func (tmdb *TMDB) GetTVEpisodeCredits(seriesID int, seasonNumber int, episodeNumber int, language *string) (*TVCreditsResponse, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	response := TVCreditsResponse{}
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(seriesID)+"/season/"+strconv.Itoa(seasonNumber)+"/episode/"+strconv.Itoa(episodeNumber)+"/credits",
		params,
		nil,
		&response,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季 第 %d 集」演员列表失败：%v", seriesID, seasonNumber, episodeNumber, err))
	}
	return &response, nil
}

type TVEpisodeImage struct {
	ID     int       `json:"id"`
	Stills []TVImage `json:"stills"`
}

// 获取属于电视剧单集的图片。
// Get the images that belong to a TV episode.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/episode/{episode_number}/images
// https://developer.themoviedb.org/reference/tv-episode-images
func (tmdb *TMDB) GetTVEpisodeImage(seriesID int, seasonNumber int, episodeNumber int, language *string) (*TVEpisodeImage, error) {
	var img TVEpisodeImage

	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(seriesID))+"/season/"+strconv.Itoa(int(seasonNumber))+"/episode/"+strconv.Itoa(int(episodeNumber))+"/images",
		params,
		nil,
		&img,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季 第 %d 集」图片失败：%v", seriesID, seasonNumber, episodeNumber, err))
	}
	return &img, nil
}
