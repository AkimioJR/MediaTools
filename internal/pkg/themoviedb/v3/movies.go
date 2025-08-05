package themoviedb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type MovieDetail struct {
	Adult               bool   `json:"adult"`
	BackdropPath        string `json:"backdrop_path"`
	BelongsToCollection string `json:"belongs_to_collection"`
	Budget              int    `json:"budget"`
	Genres              []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	Homepage            string  `json:"homepage"`
	ID                  int     `json:"id"`
	ImdbID              string  `json:"imdb_id"`
	OriginalLanguage    string  `json:"original_language"`
	OriginalTitle       string  `json:"original_title"`
	Overview            string  `json:"overview"`
	Popularity          float64 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
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
	ReleaseDate     string `json:"release_date"`
	Revenue         int    `json:"revenue"`
	Runtime         int    `json:"runtime"`
	SpokenLanguages []struct {
		EnglishName string `json:"english_name"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
	} `json:"spoken_languages"`
	Status      string  `json:"status"`
	Tagline     string  `json:"tagline"`
	Title       string  `json:"title"`
	Video       bool    `json:"video"`
	VoteAverage float64 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
}

// 通过ID获取一部电影的顶级详情。
// Get the top level details of a movie by ID.
// https://api.themoviedb.org/3/movie/{movie_id}
// https://developer.themoviedb.org/reference/movie-details
func (tmdb *TMDB) GetMovieDetail(movieID int, language *string) (*MovieDetail, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	detail := MovieDetail{}
	err := tmdb.DoRequest(http.MethodGet, "/movie/"+strconv.Itoa(movieID), params, nil, &detail)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影详情失败：%v", err))
	}
	return &detail, nil
}

type MovieAlternativeTitle struct {
	ID     int `json:"id"`
	Titles []struct {
		Iso31661 string `json:"iso_3166_1"`
		Title    string `json:"title"`
		Type     string `json:"type"`
	} `json:"titles"`
}

// 获取一部电影的其他标题。
// Get the alternative titles for a movie.
// https://api.themoviedb.org/3/movie/{movie_id}/alternative_titles
// https://developer.themoviedb.org/reference/movie-alternative-titles
//
// country 可选，指定国家(指定一个 ISO-3166-1 值来筛选结果)
func (tmdb *TMDB) GetMovieAlternativeTitle(movieID int, country *string) (*MovieAlternativeTitle, error) {
	params := url.Values{}
	if country != nil {
		params.Set("country", *country)
	}

	var resp MovieAlternativeTitle

	err := tmdb.DoRequest(http.MethodGet, "/movie/"+strconv.Itoa(movieID)+"/alternative_titles", params, nil, &resp)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影「%d」别名失败：%v", movieID, err))
	}
	return &resp, nil
}

type MovieCredit struct {
	ID   int `json:"id"`
	Cast []struct {
		Adult              bool    `json:"adult"`
		Gender             int     `json:"gender"`
		ID                 int     `json:"id"`
		KnownForDepartment string  `json:"known_for_department"`
		Name               string  `json:"name"`
		OriginalName       string  `json:"original_name"`
		Popularity         float64 `json:"popularity"`
		ProfilePath        string  `json:"profile_path"`
		CastID             int     `json:"cast_id"`
		Character          string  `json:"character"`
		CreditID           string  `json:"credit_id"`
		Order              int     `json:"order"`
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
		CreditID           string  `json:"credit_id"`
		Department         string  `json:"department"`
		Job                string  `json:"job"`
	} `json:"crew"`
}

// 获取一部电影的演员列表和工作人员列表。
// Get the cast and crew for a movie by its ID.
// https://api.themoviedb.org/3/movie/{movie_id}/credits
// https://developer.themoviedb.org/reference/movie-credits
func (tmdb *TMDB) GetMovieCredit(movieID int, language *string) (*MovieCredit, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	var response MovieCredit
	err := tmdb.DoRequest(http.MethodGet, "/movie/"+strconv.Itoa(movieID)+"/credits", params, nil, &response)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影「%d」演员列表失败：%v", movieID, err))
	}
	return &response, nil
}

type MovieExternalID struct {
	ID         int    `json:"id"`
	ImdbID     string `json:"imdb_id"`
	WikidataID string `json:"wikidata_id"`
	FacebookID string `json:"facebook_id"`
	TwitterID  string `json:"twitter_id"`
}

// 获取一部电影的外部ID。
// Get the external IDs for a movie by its ID.
// https://api.themoviedb.org/3/movie/{movie_id}/external_ids
// https://developer.themoviedb.org/reference/movie-external-ids
func (tmdb *TMDB) GetMovieExternalID(movieID int) (*MovieExternalID, error) {
	var resp MovieExternalID
	err := tmdb.DoRequest(
		http.MethodGet,
		"/movie/"+strconv.Itoa(movieID)+"/external_ids",
		url.Values{},
		nil,
		&resp,
	)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影「%d」外部ID失败：%v", movieID, err))
	}
	return &resp, nil
}

type MovieImage struct {
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

// 获取属于某部电影的图片。
// Get the images that belong to a movie.
// https://api.themoviedb.org/3/movie/{movie_id}/images
// https://developer.themoviedb.org/reference/movie-images
func (tmdb *TMDB) GetMovieImage(movieID int, language *string, IncludeImageLanguage *string) (*MovieImage, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	if IncludeImageLanguage != nil {
		params.Set("include_image_language", *IncludeImageLanguage)
	} else {
		params.Set("include_image_language", tmdb.imageLanguage)
	}

	img := MovieImage{}
	err := tmdb.DoRequest(http.MethodGet, "/movie/"+strconv.Itoa(movieID)+"/images", params, nil, &img)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影「%d」图片失败：%v", movieID, err))
	}
	return &img, nil
}

type MovieKeyword struct {
	ID       int `json:"id"`
	Keywords []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"keywords"`
}

// 获取一部电影的关键词列表。
// Get the keywords for a movie by its ID.
// https://api.themoviedb.org/3/movie/{movie_id}/keywords
// https://developer.themoviedb.org/reference/movie-keywords
func (tmdb *TMDB) GetMovieKeyword(movieID int) (*MovieKeyword, error) {
	keyword := MovieKeyword{}
	err := tmdb.DoRequest(http.MethodGet, "/movie/"+strconv.Itoa(int(movieID))+"/keywords", url.Values{}, nil, &keyword)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影「%d」关键词失败：%v", movieID, err))
	}
	if keyword.ID == 0 {
		return nil, NewTMDBError(nil, fmt.Sprintf("电影「%d」不存在或没有关键词", movieID))
	}
	return &keyword, nil
}

type MovieTranslation struct {
	ID           int `json:"id"`
	Translations []struct {
		Iso31661    string `json:"iso_3166_1"`
		Iso6391     string `json:"iso_639_1"`
		Name        string `json:"name"`
		EnglishName string `json:"english_name"`
		Data        struct {
			Homepage string `json:"homepage"`
			Overview string `json:"overview"`
			Runtime  int    `json:"runtime"`
			Tagline  string `json:"tagline"`
			Title    string `json:"title"`
		} `json:"data"`
	} `json:"translations"`
}

// 获取一部电影的翻译内容。
// Get the translations for a movie.
// https://api.themoviedb.org/3/movie/{movie_id}/translations
// https://developer.themoviedb.org/reference/movie-translations
func (tmdb *TMDB) GetMovieTranslation(movieID int) (*MovieTranslation, error) {
	params := url.Values{}
	var resp MovieTranslation
	err := tmdb.DoRequest(http.MethodGet, "/movie/"+strconv.Itoa(movieID)+"/translations", params, nil, &resp)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电影「%d」翻译失败：%v", movieID, err))
	}
	return &resp, nil
}
