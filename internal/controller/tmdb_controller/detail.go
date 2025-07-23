package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/pkg/themoviedb/v3"
	"MediaTools/internal/schemas"
	"fmt"

	"github.com/sirupsen/logrus"
)

// 搜索tmdb中所有的标题和译名
func getNames(tmdbID int, mtype meta.MediaType) ([]string, error) {

	var (
		names        []string
		titles       []themoviedb.Title
		translations []themoviedb.Translation
		err          error
	)

	switch mtype {
	case meta.MediaTypeMovie:
		titles, err = client.GetMovieAlternativeTitles(tmdbID, nil)
		if err != nil {
			return nil, fmt.Errorf("获取电影「%d」的其他标题失败: %v", tmdbID, err)
		}
		translations, err = client.GetMovieTranslations(tmdbID)
		if err != nil {
			return nil, fmt.Errorf("获取电影「%d」的翻译失败: %v", tmdbID, err)
		}
	case meta.MediaTypeTV:
		titles, err = client.GetTVSeriesAlternativeTitles(tmdbID, nil)
		if err != nil {
			return nil, fmt.Errorf("获取电视剧「%d」的其他标题失败: %v", tmdbID, err)
		}
		translations, err = client.GetTVSeriesTranslations(tmdbID)
		if err != nil {
			return nil, fmt.Errorf("获取电视剧「%d」的翻译失败: %v", tmdbID, err)
		}
	default:
		return nil, fmt.Errorf("不支持的媒体类型: 「%s」", mtype.String())
	}

	for _, title := range titles {
		names = append(names, title.Title)
	}
	for _, translation := range translations {
		names = append(names, translation.Data.Name)
	}
	return names, nil
}

func getMovieDetail(tmdbID int) (*schemas.MediaInfo, error) {
	logrus.Infof("获取电影详情，TMDB ID: %d", tmdbID)

	detail, err := client.GetMovieDetails(tmdbID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电影详情失败，TMDB ID: %d, 错误: %v", tmdbID, err)
	}

	var mediaInfo schemas.MediaInfo
	mediaInfo.TMDBID = tmdbID
	mediaInfo.MediaType = meta.MediaTypeMovie
	mediaInfo.TMDBInfo = schemas.TMDBInfo{
		MovieInfo: detail,
	}
	return &mediaInfo, nil
}

func getTVDetail(tmdbID int) (*schemas.MediaInfo, error) {
	logrus.Infof("获取电视剧详情，TMDB ID: %d", tmdbID)

	detail, err := client.GetTVSeriesDetails(tmdbID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电视剧详情失败，TMDB ID: %d, 错误: %v", tmdbID, err)
	}

	var mediaInfo schemas.MediaInfo
	mediaInfo.TMDBID = tmdbID
	mediaInfo.MediaType = meta.MediaTypeTV
	mediaInfo.TMDBInfo = schemas.TMDBInfo{
		TVInfo: schemas.TMDBTVInfo{
			SeriesInfo: detail,
		},
	}
	return &mediaInfo, nil
}
