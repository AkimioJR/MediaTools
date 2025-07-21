package themoviedb

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TVSeasonDetailResponse struct {
	Keyword
	// _ID        string      `json:"_id"` // 文档最顶层（整个剧集对象）ID
	AirDate      string      `json:"air_date"`      // 首播日期
	Episodes     []TVEpisode `json:"episodes"`      // 季内剧集列表
	Crew         []Crew      `json:"crew"`          // 工作人员列表
	GuestStars   []GuestStar `json:"guest_stars"`   // 特邀演员列表
	PosterPath   string      `json:"poster_path"`   // 海报图片路径
	SeasonNumber uint64      `json:"season_number"` // 季数
	VoteAverage  float64     `json:"vote_average"`  // 平均评分
}

// 查询一个电视剧季的详细信息。
// Query the details of a TV season.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}
// https://developer.themoviedb.org/reference/tv-season-details
func (tmdb *TMDB) GetTVSeasonDetail(seriesID uint64, seasonNumber uint64, language *string) (*TVSeasonDetailResponse, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	resp := TVSeasonDetailResponse{}
	err := tmdb.DoRequest(http.MethodGet, "/tv/"+strconv.Itoa(int(seriesID))+"/season/"+strconv.Itoa(int(seasonNumber)), params, nil, &resp)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季」详情失败：%v", seriesID, seasonNumber, err))
	}
	return &resp, nil
}

type TVSeasonImagee struct {
	Vote
	AspectRatio float64 `json:"aspect_ratio"` // 图片宽高比
	Height      uint64  `json:"height"`       // 图片高度
	ISO         string  `json:"iso_639_1"`    // ISO 639-1 语言代码
	FilePath    string  `json:"file_path"`    // 图片文件路径
	Width       uint64  `json:"width"`        // 图片宽度
}
type TVSeasonImageResponse struct {
	ID      uint64           `json:"id"`      // 电视剧ID
	Posters []TVSeasonImagee `json:"posters"` // 海报图片列表
}

// 获取属于某一电视剧季的图片。
// Get the images that belong to a TV season.
// https://api.themoviedb.org/3/tv/{series_id}/season/{season_number}/images
// https://developer.themoviedb.org/reference/tv-season-images
func (tmdb *TMDB) GetTVSeasonImage(series_id uint64, season_number uint64, language *string) (*TVSeasonImageResponse, error) {
	params := url.Values{}
	if language != nil {
		params.Set("language", *language)
	} else {
		params.Set("language", tmdb.language)
	}

	resp := TVSeasonImageResponse{}
	err := tmdb.DoRequest(http.MethodGet, "/tv/"+strconv.Itoa(int(series_id))+"/season/"+strconv.Itoa(int(season_number))+"/images", params, nil, &resp)
	if err != nil {
		return nil, NewTMDBError(err, fmt.Sprintf("获取电视剧「%d 第 %d 季」图片失败：%v", series_id, season_number, err))
	}
	return &resp, nil
}
