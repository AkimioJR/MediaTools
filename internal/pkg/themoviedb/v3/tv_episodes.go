package themoviedb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TVEpisodeDetail struct {
	AirDate string `json:"air_date"`
	Crew    []struct {
		Job                string  `json:"job"`
		Department         string  `json:"department"`
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
	EpisodeNumber int    `json:"episode_number"`
	EpisodeType   string `json:"episode_type"`
	GuestStars    []struct {
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
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ID             int     `json:"id"`
	ProductionCode string  `json:"production_code"`
	Runtime        int     `json:"runtime"`
	SeasonNumber   int     `json:"season_number"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
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

type TVEpisodeCredit struct {
	Cast []struct {
		Adult              bool    `json:"adult"`
		Gender             int     `json:"gender"`
		ID                 int     `json:"id"`
		KnownForDepartment string  `json:"known_for_department"`
		Name               string  `json:"name"`
		OriginalName       string  `json:"original_name"`
		Popularity         float64 `json:"popularity"`
		ProfilePath        string  `json:"profile_path"`
		Character          string  `json:"character"`
		CreditID           string  `json:"credit_id"`
		Order              int     `json:"order"`
	} `json:"cast"`
	Crew []struct {
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
	ID int `json:"id"`
}

// 获取一部电视剧单集的演员列表和工作人员列表。
// Get the cast and crew for a TV episode by its ID.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/episode/{episode_number}/credits
// https://developer.themoviedb.org/reference/tv-episode-credits
func (tmdb *TMDB) GetTVEpisodeCredit(seriesID int, seasonNumber int, episodeNumber int, language *string) (*TVEpisodeCredit, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	var resp TVEpisodeCredit
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(seriesID)+"/season/"+strconv.Itoa(seasonNumber)+"/episode/"+strconv.Itoa(episodeNumber)+"/credits",
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季 第 %d 集」演员列表失败：%v", seriesID, seasonNumber, episodeNumber, err))
	}
	return &resp, nil
}

type TVEpisodeExternalID struct {
	ID          int    `json:"id"`
	ImdbID      string `json:"imdb_id"`
	FreebaseMid string `json:"freebase_mid"`
	FreebaseID  string `json:"freebase_id"`
	TvdbID      int    `json:"tvdb_id"`
	TvrageID    int    `json:"tvrage_id"`
	WikidataID  string `json:"wikidata_id"`
}

// 获取已添加到一集电视剧中的外部 ID 列表。
// Get a list of external IDs that have been added to a TV episode.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/episode/{episode_number}/external_ids
// https://developer.themoviedb.org/reference/tv-episode-external-ids
func (tmdb *TMDB) GetTVEpisodeExternalID(seriesID int, seasonNumber int, episodeNumber int) (*TVEpisodeExternalID, error) {
	var resp TVEpisodeExternalID

	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(seriesID)+"/season/"+strconv.Itoa(seasonNumber)+"/episode/"+strconv.Itoa(episodeNumber)+"/external_ids",
		url.Values{},
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季 第 %d 集」外部 ID 列表失败：%v", seriesID, seasonNumber, episodeNumber, err))
	}
	return &resp, nil
}

type TVEpisodeImage struct {
	ID     int `json:"id"`
	Stills []struct {
		AspectRatio float64 `json:"aspect_ratio"`
		Height      int     `json:"height"`
		Iso6391     string  `json:"iso_639_1"`
		FilePath    string  `json:"file_path"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
		Width       int     `json:"width"`
	} `json:"stills"`
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
