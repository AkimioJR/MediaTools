package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"

	"github.com/sirupsen/logrus"
)

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

func GetMovieDetail(movieID int) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电影（TMDB ID: %d）详情", movieID)

	detail, err := client.GetMovieDetail(movieID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电影（TMDB ID: %d）详情失败: %v", movieID, err)
	}

	var mediaInfo schemas.MediaInfo
	mediaInfo.TMDBID = movieID
	mediaInfo.MediaType = meta.MediaTypeMovie
	mediaInfo.TMDBInfo = schemas.TMDBInfo{
		MovieInfo: detail,
	}
	externalID, err := client.GetMovieExternalID(movieID)
	if err != nil {
		return nil, fmt.Errorf("获取电影（TMDB ID: %d）外部ID失败: %v", movieID, err)
	}
	mediaInfo.IMDBID = externalID.ImdbID
	return &mediaInfo, nil
}

func GetTVSerieDetail(seriesID int) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）详情", seriesID)
	detail, err := client.GetTVSerieDetail(seriesID, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电视剧（TMDB ID: %d）详情失败: %v", seriesID, err)
	}

	var mediaInfo schemas.MediaInfo
	mediaInfo.TMDBID = seriesID
	mediaInfo.MediaType = meta.MediaTypeTV
	mediaInfo.TMDBInfo = schemas.TMDBInfo{
		TVInfo: schemas.TMDBTVInfo{
			SerieInfo:     detail,
			SeasonNumber:  -1, // 默认值 -1，表示未指定季数
			EpisodeNumber: -1, // 默认值 -1，表示未指定集数
		},
	}

	externalID, err := client.GetTVSerieExternalID(seriesID)
	if err != nil {
		return nil, fmt.Errorf("获取电视剧（TMDB ID: %d）外部ID失败: %v", seriesID, err)
	}
	mediaInfo.IMDBID = externalID.ImdbID
	mediaInfo.TVDBID = externalID.TvdbID

	return &mediaInfo, nil
}

func GetTVSeasonDetail(seriesID int, seasonNumber int) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）第 %d 季详情", seriesID, seasonNumber)
	seasonDetail, err := client.GetTVSeasonDetail(seriesID, seasonNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电视剧（TMDB ID: %d）第 %d 季详情失败: %v", seriesID, seasonNumber, err)
	}

	var mediaInfo schemas.MediaInfo
	mediaInfo.TMDBID = seriesID
	mediaInfo.MediaType = meta.MediaTypeTV
	mediaInfo.TMDBInfo = schemas.TMDBInfo{
		TVInfo: schemas.TMDBTVInfo{
			SeasonInfo:    seasonDetail,
			SeasonNumber:  seasonNumber,
			EpisodeNumber: -1, // 默认值 -1，表示未指定集数
		},
	}
	externalID, err := client.GetTVSeasonExternalID(seriesID, seasonNumber)
	if err != nil {
		return nil, fmt.Errorf("获取电视剧（TMDB ID: %d）第 %d 季外部ID失败, 错误: %v", seriesID, seasonNumber, err)
	}
	mediaInfo.TVDBID = externalID.TvdbID

	return &mediaInfo, nil
}

func GetTVEpisodeDetail(seriesID int, seasonNumber int, episodeNumber int) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）第 %d 季第 %d 集详情", seriesID, seasonNumber, episodeNumber)
	episodeDetail, err := client.GetTVEpisodeDetail(seriesID, seasonNumber, episodeNumber, nil)
	if err != nil {
		return nil, fmt.Errorf("获取电视剧（TMDB ID: %d）第 %d 季第 %d 集详情失败，错误: %v", seriesID, seasonNumber, episodeNumber, err)
	}

	var mediaInfo schemas.MediaInfo
	mediaInfo.TMDBID = seriesID
	mediaInfo.MediaType = meta.MediaTypeTV
	mediaInfo.TMDBInfo = schemas.TMDBInfo{
		TVInfo: schemas.TMDBTVInfo{
			EpisodeInfo:   episodeDetail,
			SeasonNumber:  seasonNumber,
			EpisodeNumber: episodeNumber,
		},
	}
	externalID, err := client.GetTVEpisodeExternalID(seriesID, seasonNumber, episodeNumber)
	if err != nil {
		return nil, fmt.Errorf("获取电视剧（TMDB ID: %d）第 %d 季第 %d 集外部ID失败，错误: %v", seriesID, seasonNumber, episodeNumber, err)
	}
	mediaInfo.IMDBID = externalID.ImdbID
	mediaInfo.TVDBID = externalID.TvdbID

	return &mediaInfo, nil
}
