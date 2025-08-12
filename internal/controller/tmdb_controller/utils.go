package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"

	"github.com/sirupsen/logrus"
)

// 给定TMDB号，查询一条媒体信息
// tmdbID TMDB ID
// mtype 媒体类型，未指定（MediaTypeUnknown）时会尝试识别
func GetInfo(tmdbID int, mtype meta.MediaType) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	if mtype == meta.MediaTypeUnknown {
		logrus.Infof("未指定 TMDB ID 「%d」的媒体类型", tmdbID)
		movieDetail, movieErr := GetMovieDetail(tmdbID)
		tvDetail, tvErr := GetTVSerieDetail(tmdbID)

		switch {
		case movieErr == nil && tvErr == nil:
			return movieDetail, fmt.Errorf("TMDB ID 「%d」同时匹配到电影和电视剧，无法识别", tmdbID)
		case movieErr == nil:
			logrus.Infof("识别为电影，TMDB ID: %d", tmdbID)
			return movieDetail, nil
		case tvErr == nil:
			logrus.Infof("识别为电视剧，TMDB ID: %d", tmdbID)
			return tvDetail, nil
		default:
			return nil, fmt.Errorf("未查询到 TMDB ID「%d」信息", tmdbID)
		}
	}

	switch mtype {
	case meta.MediaTypeMovie:
		return GetMovieDetail(tmdbID)
	case meta.MediaTypeTV:
		return GetTVSerieDetail(tmdbID)
	default:
		return nil, fmt.Errorf("不支持的媒体类型: 「%s」", mtype)
	}
}

// 搜索tmdb中所有的标题和译名
func getNames(tmdbID int, mtype meta.MediaType) ([]string, error) {
	lock.RLock()
	defer lock.RUnlock()

	var names []string

	switch mtype {
	case meta.MediaTypeMovie:
		titleResp, err := client.GetMovieAlternativeTitle(tmdbID, nil)
		if err != nil {
			return nil, fmt.Errorf("获取电影「%d」的其他标题失败: %v", tmdbID, err)
		}
		translationResp, err := client.GetMovieTranslation(tmdbID)
		if err != nil {
			return nil, fmt.Errorf("获取电影「%d」的翻译失败: %v", tmdbID, err)
		}
		for _, title := range titleResp.Titles {
			names = append(names, title.Title)
		}
		for _, translation := range translationResp.Translations {
			names = append(names, translation.Data.Title)
		}
	case meta.MediaTypeTV:
		titleResp, err := client.GetTVSerieAlternativeTitle(tmdbID, nil)
		if err != nil {
			return nil, fmt.Errorf("获取电视剧「%d」的其他标题失败: %v", tmdbID, err)
		}
		translationResp, err := client.GetTVSerieTranslation(tmdbID)
		if err != nil {
			return nil, fmt.Errorf("获取电视剧「%d」的翻译失败: %v", tmdbID, err)
		}
		for _, title := range titleResp.Results {
			names = append(names, title.Title)
		}
		for _, translation := range translationResp.Translations {
			names = append(names, translation.Data.Name)
		}
	default:
		return nil, fmt.Errorf("不支持的媒体类型: 「%s」", mtype.String())
	}
	return names, nil
}

func parseType(s string) meta.MediaType {
	switch string(s) {
	case "movie":
		return meta.MediaTypeMovie
	case "tv":
		return meta.MediaTypeTV
	default:
		return meta.MediaTypeUnknown
	}
}
