package tmdb_controller

import (
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"

	"github.com/sirupsen/logrus"
)

// RecognizeMedia 识别媒体信息
// videoMeta 识别的元数据
// 返回识别后的媒体信息
func RecognizeMedia(videoMeta *meta.VideoMeta) (*schemas.MediaInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	// 如果视频元数据中包含 TMDB ID，则直接查询
	if videoMeta.TMDBID > 0 {
		return GetInfo(videoMeta.TMDBID, videoMeta.MediaType)
	}

	// 如果没有 TMDB ID，则尝试识别媒体名称
	var fn func(string) (*schemas.MediaInfo, error)
	switch videoMeta.MediaType {
	case meta.MediaTypeMovie:
		fn = SearchMovie
	case meta.MediaTypeTV:
		fn = SearchTV
	case meta.MediaTypeUnknown:
		fn = SearchMulti
	}
	for _, name := range videoMeta.GetTitles() {
		info, err := fn(name)
		if err != nil {
			logrus.Warningf("识别「%s」媒体信息失败: %v", name, err)
			continue // 如果搜索失败，尝试下一个名称
		}
		return info, nil // 找到匹配的媒体信息
	}
	return nil, fmt.Errorf("未能 %v 识别媒体信息，可能是名称不匹配或 TMDB 中没有相关数据", videoMeta.GetTitles())
}

// RecognizeAndEnrichMedia 识别媒体信息并，如果是电视剧类型，还会补充季和集的详细信息（如果有对应信息）
// videoMeta 识别的元数据
// 返回识别后的媒体信息，如果是电视剧类型，还会补充季和集的详细信息
func RecognizeAndEnrichMedia(videoMeta *meta.VideoMeta) (*schemas.MediaInfo, error) {
	info, err := RecognizeMedia(videoMeta)
	if err != nil {
		return nil, fmt.Errorf("识别媒体信息失败: %v", err)
	}
	if info.MediaType == meta.MediaTypeTV {
		if videoMeta.Season != -1 {
			seasonDetail, err := GetTVSeasonDetail(info.TMDBID, videoMeta.Season)
			if err != nil {
				return nil, fmt.Errorf("获取电视剧季信息失败: %v", err)
			}
			info.TMDBInfo.TVInfo.SeasonInfo = seasonDetail.TMDBInfo.TVInfo.SeasonInfo
		}
		if videoMeta.Episode != -1 {
			episodeDetail, err := GetTVEpisodeDetail(info.TMDBID, videoMeta.Season, videoMeta.Episode)
			if err != nil {
				return nil, fmt.Errorf("获取电视剧集信息失败: %v", err)
			}
			info.TMDBInfo.TVInfo.EpisodeInfo = episodeDetail.TMDBInfo.TVInfo.EpisodeInfo
		}
	}
	return info, nil
}
