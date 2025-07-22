package themoviedb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TVEpisodeDetail struct {
	Keyword
	Vote
	AirDate        string      `json:"air_date"`        // 首播日期
	Crew           []Crew      `json:"crew"`            // 工作人员列表
	EpisodeNumber  uint64      `json:"episode_number"`  // 集数
	GuestStars     []GuestStar `json:"guest_stars"`     // 特邀演员列表
	Overview       string      `json:"overview"`        // 概述
	ProductionCode string      `json:"production_code"` // 制作代码
	Runtime        uint64      `json:"runtime"`         // 时长
	SeasonNumber   uint64      `json:"season_number"`   // 季数
	StillPath      string      `json:"still_path"`      // 静态图片路径
}

// 查询电视剧单集的详细信息。
// Query the details of a TV episode.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/episode/{episode_number}
// https://developer.themoviedb.org/reference/tv-episode-details
func (tmdb *TMDB) GetTVEpisodeDetail(seriesID uint64, seasonNumber uint64, episodeNumber uint64, language *string) (*TVEpisodeDetail, error) {
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

type TVEpisodeImage struct {
	ID     uint64    `json:"id"`
	Stills []TVImage `json:"stills"`
}

// 获取属于电视剧单集的图片。
// Get the images that belong to a TV episode.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/episode/{episode_number}/images
// https://developer.themoviedb.org/reference/tv-episode-images
func (tmdb *TMDB) GetTVEpisodeImage(seriesID uint64, seasonNumber uint64, episodeNumber uint64, language *string) (*TVEpisodeImage, error) {
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
