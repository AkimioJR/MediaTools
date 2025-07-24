package themoviedb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TVSeasonDetail struct {
	ID       string `json:"_id"`
	AirDate  string `json:"air_date"`
	Episodes []struct {
		AirDate        string  `json:"air_date"`
		EpisodeNumber  int     `json:"episode_number"`
		ID             int     `json:"id"`
		Name           string  `json:"name"`
		Overview       string  `json:"overview"`
		ProductionCode string  `json:"production_code"`
		Runtime        int     `json:"runtime"`
		SeasonNumber   int     `json:"season_number"`
		ShowID         int     `json:"show_id"`
		StillPath      string  `json:"still_path"`
		VoteAverage    float64 `json:"vote_average"`
		VoteCount      int     `json:"vote_count"`
		Crew           []struct {
			Department         string  `json:"department"`
			Job                string  `json:"job"`
			CreditID           string  `json:"credit_id"`
			Adult              bool    `json:"adult"`
			Gender             int     `json:"gender"`
			ID                 int     `json:"id"`
			KnownForDepartment string  `json:"known_for_department"`
			Name               string  `json:"name"`
			OriginalName       string  `json:"original_name"`
			Popularity         float64 `json:"popularity"`
			ProfilePath        string  `json:"profile_path"`
		} `json:"crew"`
		GuestStars []struct {
			Character          string  `json:"character"`
			CreditID           string  `json:"credit_id"`
			Order              int     `json:"order"`
			Adult              bool    `json:"adult"`
			Gender             int     `json:"gender"`
			ID                 int     `json:"id"`
			KnownForDepartment string  `json:"known_for_department"`
			Name               string  `json:"name"`
			OriginalName       string  `json:"original_name"`
			Popularity         float64 `json:"popularity"`
			ProfilePath        string  `json:"profile_path"`
		} `json:"guest_stars"`
	} `json:"episodes"`
	Name         string  `json:"name"`
	Overview     string  `json:"overview"`
	ID0          int     `json:"id"`
	PosterPath   string  `json:"poster_path"`
	SeasonNumber int     `json:"season_number"`
	VoteAverage  float64 `json:"vote_average"`
}

// 查询一个电视剧季的详细信息。
// Query the details of a TV season.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}
// https://developer.themoviedb.org/reference/tv-season-details
func (tmdb *TMDB) GetTVSeasonDetail(seriesID int, seasonNumber int, language *string) (*TVSeasonDetail, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	var resp TVSeasonDetail
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(seriesID))+"/season/"+strconv.Itoa(int(seasonNumber)),
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季」详情失败：%v", seriesID, seasonNumber, err))
	}
	return &resp, nil
}

type TVSeasonCredit struct {
	Cast []struct {
		Adult              bool    `json:"adult"`
		Gender             int     `json:"gender"`
		ID                 int     `json:"id"`
		KnownForDepartment string  `json:"known_for_department"`
		Name               string  `json:"name"`
		OriginalName       string  `json:"original_name"`
		Popularity         float64 `json:"popularity"`
		ProfilePath        string  `json:"profile_path"`
		Roles              []struct {
			CreditID     string `json:"credit_id"`
			Character    string `json:"character"`
			EpisodeCount int    `json:"episode_count"`
		} `json:"roles"`
		TotalEpisodeCount int `json:"total_episode_count"`
		Order             int `json:"order"`
	} `json:"cast"`
	Crew []struct {
		Adult              bool        `json:"adult"`
		Gender             int         `json:"gender"`
		ID                 int         `json:"id"`
		KnownForDepartment string      `json:"known_for_department"`
		Name               string      `json:"name"`
		OriginalName       string      `json:"original_name"`
		Popularity         float64     `json:"popularity"`
		ProfilePath        interface{} `json:"profile_path"`
		Jobs               []struct {
			CreditID     string `json:"credit_id"`
			Job          string `json:"job"`
			EpisodeCount int    `json:"episode_count"`
		} `json:"jobs"`
		Department        string `json:"department"`
		TotalEpisodeCount int    `json:"total_episode_count"`
	} `json:"crew"`
	ID int `json:"id"`
}

// 获取一部电视剧季的演员列表和工作人员列表。
// Get the cast and crew for a TV season by its ID.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/credits
// https://developer.themoviedb.org/reference/tv-season-credits
func (tmdb *TMDB) GetTVSeasonCredit(seriesID int, seasonNumber int, language *string) (*TVSeasonCredit, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	var resp TVSeasonCredit
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(seriesID)+"/season/"+strconv.Itoa(seasonNumber)+"/credits",
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季」演员列表失败：%v", seriesID, seasonNumber, err))
	}
	return &resp, nil
}

type TVSeasonImage struct {
	Backdrops []struct {
		AspectRatio float64     `json:"aspect_ratio"`
		Height      int         `json:"height"`
		Iso6391     interface{} `json:"iso_639_1"`
		FilePath    string      `json:"file_path"`
		VoteAverage float64     `json:"vote_average"`
		VoteCount   int         `json:"vote_count"`
		Width       int         `json:"width"`
	} `json:"backdrops"`
	ID    int `json:"id"`
	Logos []struct {
		AspectRatio float64 `json:"aspect_ratio"`
		Height      int     `json:"height"`
		Iso6391     string  `json:"iso_639_1"`
		FilePath    string  `json:"file_path"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
		Width       int     `json:"width"`
	} `json:"logos"`
	Posters []struct {
		AspectRatio float64 `json:"aspect_ratio"`
		Height      int     `json:"height"`
		Iso6391     string  `json:"iso_639_1"`
		FilePath    string  `json:"file_path"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
		Width       int     `json:"width"`
	} `json:"posters"`
}

// 获取属于某一电视剧季的图片。
// Get the images that belong to a TV season.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/images
// https://developer.themoviedb.org/reference/tv-season-images
func (tmdb *TMDB) GetTVSeasonImage(series_id int, season_number int, language *string) (*TVSeasonImage, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	var resp TVSeasonImage
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(series_id))+"/season/"+strconv.Itoa(int(season_number))+"/images",
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季」图片失败：%v", series_id, season_number, err))
	}
	return &resp, nil
}
