package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"

	"github.com/sirupsen/logrus"
)

// 给定TMDB号，查询一条媒体信息
func GetInfo(tmdbID uint64, mtype *meta.MediaType) (*schemas.MediaInfo, error) {
	if mtype == nil || *mtype == meta.MediaTypeUnknown {
		logrus.Infof("未指定 TMDB ID 「%d」的媒体类型", tmdbID)
		movieDetail, movieErr := getMovieDetail(tmdbID)
		tvDetail, tvErr := getTVDetail(tmdbID)

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
		return getMovieDetail(tmdbID)
	case meta.MediaTypeTV:
		return getTVDetail(tmdbID)
	default:
		return nil, fmt.Errorf("不支持的媒体类型: 「%s」", mtype)
	}
}

// RecognizeMedia 识别媒体信息
// videoMeta 识别的元数据
// mtype 媒体类型
// tmdbID TMDB ID
// 返回识别后的媒体信息
func RecognizeMedia(videoMeta *meta.VideoMeta, mtype *meta.MediaType, tmdbID *uint64) (*schemas.MediaInfo, error) {
	if tmdbID == nil && videoMeta == nil {
		return nil, fmt.Errorf("没有提供 TMDB ID 或元数据，无法识别媒体信息")
	}

	// if metaVideo != nil {
	// 	if mtype != nil && *mtype != meta.MediaTypeUnknown && metaVideo.GetType() != meta.MediaTypeUnknown {

	// 	}
	// }

	if tmdbID != nil { // 优先根据 TMDB ID 获取媒体信息
		return GetInfo(*tmdbID, mtype)
	}

	var titles []string
	if videoMeta.CNTitle != "" {
		titles = append(titles, videoMeta.CNTitle)
	}
	if videoMeta.ENTitle != "" {
		titles = append(titles, videoMeta.ENTitle)
	}
	for _, title := range titles {
		if videoMeta.BeginSeason == nil {
			logrus.Infof("正在识别「%s」...", title)
		} else {
			logrus.Infof("正在识别「%s（第 %d 季）」...", title, *videoMeta.BeginSeason)
		}

		switch {
		case mtype == nil, *mtype == meta.MediaTypeUnknown:
			info, err := MatchMulti(title)
			if err == nil {
				logrus.Infof("识别「%s」媒体信息成功（多媒体匹配）", title)
				return info, nil
			}
			logrus.Warningf("识别「%s」媒体信息失败: %v", title, err)

		case *mtype == meta.MediaTypeTV:
			var year *int
			if videoMeta.Year > 0 {
				intYear := int(videoMeta.Year)
				year = &intYear
			}
			info, err := Match(title, *mtype, nil, year, videoMeta.BeginSeason)
			if err == nil {
				logrus.Infof("识别「%s」电视剧信息成功", title)
				return info, nil
			}

			logrus.Warningf("识别「%s」电视剧信息失败: %v", title, err)
			info, err = MatchMulti(title) // 去掉年份和类型再查一次
			if err == nil {
				logrus.Infof("识别「%s」电视剧信息成功（多媒体匹配）", title)
				return info, nil
			}
			logrus.Warningf("识别「%s」电视剧信息失败（多媒体匹配）: %v", title, err)
		}

	}
	return nil, fmt.Errorf("无法识别媒体信息，未找到匹配结果")

}
