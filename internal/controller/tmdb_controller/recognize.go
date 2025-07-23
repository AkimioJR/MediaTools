package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"

	"github.com/sirupsen/logrus"
)

// 给定TMDB号，查询一条媒体信息
func GetInfo(tmdbID int, mtype *meta.MediaType) (*schemas.MediaInfo, error) {
	if mtype == nil || *mtype == meta.MediaTypeUnknown {
		logrus.Infof("未指定 TMDB ID 「%d」的媒体类型", tmdbID)
		movieDetail, movieErr := GetMovieDetail(tmdbID)
		tvDetail, tvErr := GetTVSeriesDetail(tmdbID)

		switch {
		case movieErr == nil && tvErr == nil:
			logrus.Warningf("TMDB ID 「%d」同时匹配到电影和电视剧，无法识别", tmdbID)
			movieDetail.MediaType = meta.MediaTypeUnknown
			movieDetail.TMDBInfo.TVInfo = tvDetail.TMDBInfo.TVInfo // 合并电视剧信息
			return movieDetail, nil
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

	switch *mtype {
	case meta.MediaTypeMovie:
		return GetMovieDetail(tmdbID)
	case meta.MediaTypeTV:
		return GetTVSeriesDetail(tmdbID)
	default:
		return nil, fmt.Errorf("不支持的媒体类型: 「%s」", mtype)
	}
}

// RecognizeMedia 识别媒体信息
// videoMeta 识别的元数据
// mtype 媒体类型
// tmdbID TMDB ID
// 返回识别后的媒体信息
func RecognizeMedia(videoMeta *meta.VideoMeta, mType *meta.MediaType, tmdbID *int) (*schemas.MediaInfo, error) {
	// 1. 优先处理直接提供 tmdbID 的情况
	if tmdbID != nil && *tmdbID > 0 {
		return GetInfo(*tmdbID, mType)
	}

	// 2. 尝试从 videoMeta 获取 tmdbID
	if videoMeta != nil && videoMeta.TMDBID > 0 {
		return GetInfo(videoMeta.TMDBID, mType)
	}

	// 3. 检查是否有有效数据
	if videoMeta == nil {
		return nil, fmt.Errorf("没有提供有效的元数据，无法识别媒体信息")
	}

	// 4. 使用指定媒体类型或元数据中的类型
	mediaType := meta.MediaTypeUnknown
	if mType != nil {
		mediaType = *mType
	} else if videoMeta.MediaType != meta.MediaTypeUnknown {
		mediaType = videoMeta.MediaType
	}

	for _, title := range videoMeta.GetTitles() {
		if videoMeta.Season == -1 {
			logrus.Infof("正在识别「%s」...", title)
		} else {
			logrus.Infof("正在识别「%s（第 %d 季）」...", title, videoMeta.Season)
		}

		var year *int
		if videoMeta.Year > 0 {
			intYear := int(videoMeta.Year)
			year = &intYear
		}
		info, err := MatchWithFallback(title, mediaType, year, year, &videoMeta.Season)
		if err == nil {
			logrus.Infof("识别「%s」媒体信息成功", title)
			return info, nil
		}
		logrus.Warningf("识别「%s」媒体信息失败: %v", title, err)
	}
	return nil, fmt.Errorf("无法识别媒体信息，未找到匹配结果")

}
