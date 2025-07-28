package themoviedb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TVSerieDetail struct {
	Adult        bool   `json:"adult"`
	BackdropPath string `json:"backdrop_path"`
	CreatedBy    []struct {
		ID          int    `json:"id"`
		CreditID    string `json:"credit_id"`
		Name        string `json:"name"`
		Gender      int    `json:"gender"`
		ProfilePath string `json:"profile_path"`
	} `json:"created_by"`
	EpisodeRunTime []int  `json:"episode_run_time"`
	FirstAirDate   string `json:"first_air_date"`
	Genres         []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage         string   `json:"homepage"`
	ID               int      `json:"id"`
	InProduction     bool     `json:"in_production"`
	Languages        []string `json:"languages"`
	LastAirDate      string   `json:"last_air_date"`
	LastEpisodeToAir struct {
		ID             int     `json:"id"`
		Name           string  `json:"name"`
		Overview       string  `json:"overview"`
		VoteAverage    float64 `json:"vote_average"`
		VoteCount      int     `json:"vote_count"`
		AirDate        string  `json:"air_date"`
		EpisodeNumber  int     `json:"episode_number"`
		ProductionCode string  `json:"production_code"`
		Runtime        int     `json:"runtime"`
		SeasonNumber   int     `json:"season_number"`
		ShowID         int     `json:"show_id"`
		StillPath      string  `json:"still_path"`
	} `json:"last_episode_to_air"`
	Name             string `json:"name"`
	NextEpisodeToAir struct {
		ID             int     `json:"id"`
		Name           string  `json:"name"`
		Overview       string  `json:"overview"`
		VoteAverage    float64 `json:"vote_average"`
		VoteCount      int     `json:"vote_count"`
		AirDate        string  `json:"air_date"`
		EpisodeNumber  int     `json:"episode_number"`
		ProductionCode string  `json:"production_code"`
		Runtime        int     `json:"runtime"`
		SeasonNumber   int     `json:"season_number"`
		ShowID         int     `json:"show_id"`
		StillPath      string  `json:"still_path"`
	} `json:"next_episode_to_air"`
	Networks []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"networks"`
	NumberOfEpisodes    int      `json:"number_of_episodes"`
	NumberOfSeasons     int      `json:"number_of_seasons"`
	OriginCountry       []string `json:"origin_country"`
	OriginalLanguage    string   `json:"original_language"`
	OriginalName        string   `json:"original_name"`
	Overview            string   `json:"overview"`
	Popularity          float64  `json:"popularity"`
	PosterPath          string   `json:"poster_path"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		LogoPath      string `json:"logo_path"`
		Name          string `json:"name"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	ProductionCountries []struct {
		Iso31661 string `json:"iso_3166_1"`
		Name     string `json:"name"`
	} `json:"production_countries"`
	Seasons []struct {
		AirDate      string  `json:"air_date"`
		EpisodeCount int     `json:"episode_count"`
		ID           int     `json:"id"`
		Name         string  `json:"name"`
		Overview     string  `json:"overview"`
		PosterPath   string  `json:"poster_path"`
		SeasonNumber int     `json:"season_number"`
		VoteAverage  float64 `json:"vote_average"`
	} `json:"seasons"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Type        string  `json:"type"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

// 获取一部电视剧的详细信息。
// Get the details of a TV show.
// https://api.themoviedb.org/3/tv/{series_id}
// https://developer.themoviedb.org/reference/tv-series-details
func (tmdb *TMDB) GetTVSerieDetail(seriesID int, language *string) (*TVSerieDetail, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	var resp TVSerieDetail
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(seriesID)),
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧详情失败：%v", err))
	}
	return &resp, nil
}

type TVSerieCredit struct {
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
		Adult              bool    `json:"adult"`
		Gender             int     `json:"gender"`
		ID                 int     `json:"id"`
		KnownForDepartment string  `json:"known_for_department"`
		Name               string  `json:"name"`
		OriginalName       string  `json:"original_name"`
		Popularity         float64 `json:"popularity"`
		ProfilePath        string  `json:"profile_path"`
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

// 获取一部电视剧最新一季的演职人员名单。
// Get the latest season credits of a TV show.
// https://api.themoviedb.org/3/tv/{series_id}/credits
// https://developer.themoviedb.org/reference/tv-series-credits
func (tmdb *TMDB) GetTVSerieCredit(seriesID int, language *string) (*TVSerieCredit, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}
	var resp TVSerieCredit
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(seriesID),
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」演员列表失败：%v", seriesID, err))
	}
	return &resp, nil
}

type TVSerieAlternativeTitle struct {
	ID      int `json:"id"`
	Results []struct {
		Iso31661 string `json:"iso_3166_1"`
		Title    string `json:"title"`
		Type     string `json:"type"`
	} `json:"results"`
}

// 获取已添加到电视剧中的其他标题。
// Get the alternative titles that have been added to a TV show.
// https://api.themoviedb.org/3/tv/{series_id}/alternative_titles
// https://developer.themoviedb.org/reference/tv-series-alternative-titles
//
// country 可选，指定国家(指定一个 ISO-3166-1 值来筛选结果)
func (tmdb *TMDB) GetTVSerieAlternativeTitle(seriesID int, country *string) (*TVSerieAlternativeTitle, error) {
	params := url.Values{}
	if country != nil {
		params.Set("country", *country)
	}

	var resp TVSerieAlternativeTitle
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(seriesID))+"/alternative_titles",
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」别名失败：%v", seriesID, err))
	}
	return &resp, nil
}

type TVSerieEpisodeGroup struct {
	Results []struct {
		Description  string `json:"description"`
		EpisodeCount int    `json:"episode_count"`
		GroupCount   int    `json:"group_count"`
		ID           string `json:"id"`
		Name         string `json:"name"`
		Network      struct {
			ID            int    `json:"id"`
			LogoPath      string `json:"logo_path"`
			Name          string `json:"name"`
			OriginCountry string `json:"origin_country"`
		} `json:"network"`
		Type int `json:"type"`
	} `json:"results"`
	ID int `json:"id"`
}

// 获取已添加到电视剧中的剧集组。
// Get the episode groups that have been added to a TV show.
// https://api.themoviedb.org/3/tv/{series_id}/episode_groups
// https://developer.themoviedb.org/reference/tv-series-episode-groups
func (tmdb *TMDB) GetTVSerieEpisodeGroup(seriesID int) (*TVSerieEpisodeGroup, error) {
	params := url.Values{}
	var resp TVSerieEpisodeGroup
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(seriesID))+"/episode_groups",
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」剧集组失败：%v", seriesID, err))
	}
	return &resp, nil
}

type TVSerieExternalID struct {
	ID          int    `json:"id"`
	ImdbID      string `json:"imdb_id"`
	FreebaseMid string `json:"freebase_mid"`
	FreebaseID  string `json:"freebase_id"`
	TvdbID      int    `json:"tvdb_id"`
	TvrageID    int    `json:"tvrage_id"`
	WikidataID  string `json:"wikidata_id"`
	FacebookID  string `json:"facebook_id"`
	InstagramID string `json:"instagram_id"`
	TwitterID   string `json:"twitter_id"`
}

// 获取已添加到一部电视剧中的外部 ID 列表。
// Get a list of external IDs that have been added to a TV show.
// https://api.themoviedb.org/3/tv/{series_id}/external_ids
// https://developer.themoviedb.org/reference/tv-series-external-ids
func (tmdb *TMDB) GetTVSerieExternalID(seriesID int) (*TVSerieExternalID, error) {
	var resp TVSerieExternalID
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(seriesID)+"/external_ids",
		url.Values{},
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」外部ID失败：%v", seriesID, err))
	}
	return &resp, nil
}

type TVSerieImage struct {
	Backdrops []struct {
		AspectRatio float64 `json:"aspect_ratio"`
		Height      int     `json:"height"`
		Iso6391     string  `json:"iso_639_1"`
		FilePath    string  `json:"file_path"`
		VoteAverage float64 `json:"vote_average"`
		VoteCount   int     `json:"vote_count"`
		Width       int     `json:"width"`
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

// 获取属于某部电视剧的图片。
// Get the images that belong to a TV series.
// https://api.themoviedb.org/3/tv/{series_id}/images
// https://developer.themoviedb.org/reference/tv-series-images
func (tmdb *TMDB) GetTVSerieImage(seriesID int, IncludeImageLanguage *string, language *string) (*TVSerieImage, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	if IncludeImageLanguage != nil {
		params.Set("include_image_language", *IncludeImageLanguage)
	}

	var resp TVSerieImage
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(seriesID))+"/images",
		params,
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」图片失败：%v", seriesID, err))
	}
	return &resp, nil
}

type TVSerieKeyword struct {
	ID      int `json:"id"`
	Results []struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"results"`
}

// 获取一部电视剧的关键词列表。
// Get the keywords for a movie by its ID.
// https://api.themoviedb.org/3/tv/{series_id}/keywords
// https://developer.themoviedb.org/reference/movie-keywords
func (tmdb *TMDB) GetTVSerieKeyword(seriesID int) (*TVSerieKeyword, error) {
	var resp TVSerieKeyword
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(seriesID))+"/keywords",
		url.Values{},
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」关键词失败：%v", seriesID, err))
	}
	return &resp, nil
}

type TVSerieTranslation struct {
	ID           int `json:"id"`
	Translations []struct {
		Iso31661    string `json:"iso_3166_1"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
		EnglishName string `json:"english_name"`
		Data        struct {
			Name     string `json:"name"`
			Overview string `json:"overview"`
			Homepage string `json:"homepage"`
			Tagline  string `json:"tagline"`
		} `json:"data"`
	} `json:"translations"`
}

// 获取已添加到电视剧中的翻译内容。
// Get the translations that have been added to a TV show.
// https://api.themoviedb.org/3/tv/{series_id}/translations
// https://developer.themoviedb.org/reference/tv-series-translations
func (tmdb *TMDB) GetTVSerieTranslation(seriesID int) (*TVSerieTranslation, error) {
	var resp TVSerieTranslation
	err := tmdb.DoRequest(
		http.MethodGet,
		"/tv/"+strconv.Itoa(int(seriesID))+"/translations",
		url.Values{},
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d」翻译失败：%v", seriesID, err))
	}
	return &resp, nil
}
