package scrape_controller

import (
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

func genMovieMetaInfo(mediaInfo *schemas.MediaInfo) *MovieMetaData {

	year, err := strconv.Atoi(mediaInfo.TMDBInfo.MovieInfo.ReleaseDate[:4])
	if err != nil {
		logrus.Warningf("获取电影「%s」发行年份失败: %v", mediaInfo.TMDBInfo.MovieInfo.Title, err)
	}

	var (
		actors    []Actor   // 演员列表
		directors []Creator // 导演列表
		writers   []Creator // 编剧列表
		credits   []Creator // 其他制作人员列表
	)
	resp, err := tmdb_controller.GetMovieCredit(mediaInfo.TMDBID, nil)
	if err != nil {
		logrus.Warningf("获取电影「%s」演员/制作人员列表失败: %v", mediaInfo.TMDBInfo.MovieInfo.Title, err)
	} else {
		for _, cast := range resp.Cast {
			actors = append(actors, Actor{
				Name:    cast.Name,
				Role:    cast.Character,
				Type:    "Actor",
				TMDBID:  strconv.Itoa(cast.ID),
				Thumb:   tmdb_controller.GetImageURL(cast.ProfilePath),
				Profile: fmt.Sprintf("https://www.themoviedb.org/person/%d", cast.ID),
			})
		}

		for _, crew := range resp.Crew {
			creator := Creator{
				Name:    crew.Name,
				TMDBID:  strconv.Itoa(crew.ID),
				Profile: fmt.Sprintf("https://www.themoviedb.org/person/%d", crew.ID),
			}
			switch crew.Job {
			case "Director":
				directors = append(directors, creator)
			case "Writer":
				writers = append(writers, creator)
			default:
				credits = append(credits, creator)
			}
		}
	}
	var genres []string
	for _, genre := range mediaInfo.TMDBInfo.MovieInfo.Genres {
		genres = append(genres, genre.Name)
	}

	var uniqueIDs []UniqueID
	if mediaInfo.TMDBID != 0 {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "tmdbid",
			Value: strconv.Itoa(mediaInfo.TMDBID),
		})
	}
	if mediaInfo.IMDBID != "" {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "imdbid",
			Value: mediaInfo.IMDBID,
		})
	}
	if mediaInfo.TVDBID != 0 {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "tvdbid",
			Value: strconv.Itoa(mediaInfo.TVDBID),
		})
	}

	data := MovieMetaData{
		Title:         mediaInfo.TMDBInfo.MovieInfo.Title,         // 电影标题
		OriginalTitle: mediaInfo.TMDBInfo.MovieInfo.OriginalTitle, // 原始标题
		Plot:          mediaInfo.TMDBInfo.MovieInfo.Overview,      // 剧情简介
		Outline:       mediaInfo.TMDBInfo.MovieInfo.Overview,      // 大纲简介
		Actors:        actors,                                     // 演员列表
		Directors:     directors,                                  // 导演列表
		Writers:       writers,                                    // 编剧列表
		Credits:       credits,                                    // 其他制作人员列表
		Genres:        genres,                                     // 电影类型
		Year:          year,                                       // 发行年份
		Premiered:     mediaInfo.TMDBInfo.MovieInfo.ReleaseDate,   // 首映日期
		Rating:        mediaInfo.TMDBInfo.MovieInfo.VoteAverage,   // 评分
		TMDBID:        mediaInfo.TMDBID,                           // TMDB ID
		UniqueIDs:     uniqueIDs,                                  // 唯一ID列表
	}

	return &data
}

func genTVSerieMetaInfo(mediaInfo *schemas.MediaInfo) *TVSeriesMetaData {
	year, err := strconv.Atoi(mediaInfo.TMDBInfo.TVInfo.SerieInfo.FirstAirDate[:4])
	if err != nil {
		logrus.Warningf("获取电视剧「%s」发行年份失败: %v", mediaInfo.TMDBInfo.TVInfo.SerieInfo.Name, err)
	}
	var actors []Actor
	resp, err := tmdb_controller.GetTVSerieCredit(mediaInfo.TMDBID, nil)
	if err != nil {
		logrus.Warningf("获取电视剧「%s」演员/制作人员列表失败: %v", mediaInfo.TMDBInfo.TVInfo.SerieInfo.Name, err)
	} else {
		for _, cast := range resp.Cast {
			for _, role := range cast.Roles {
				actors = append(actors, Actor{
					Name:    cast.Name,
					Role:    role.Character,
					Type:    "Actor",
					TMDBID:  strconv.Itoa(cast.ID),
					Thumb:   tmdb_controller.GetImageURL(cast.ProfilePath),
					Profile: fmt.Sprintf("https://www.themoviedb.org/person/%d", cast.ID),
				})
			}
		}
	}
	var genres []string
	for _, genre := range mediaInfo.TMDBInfo.TVInfo.SerieInfo.Genres {
		genres = append(genres, genre.Name)
	}

	var uniqueIDs []UniqueID
	if mediaInfo.TMDBID != 0 {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "tmdbid",
			Value: strconv.Itoa(mediaInfo.TMDBID),
		})
	}
	if mediaInfo.IMDBID != "" {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "imdbid",
			Value: mediaInfo.IMDBID,
		})
	}
	if mediaInfo.TVDBID != 0 {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "tvdbid",
			Value: strconv.Itoa(mediaInfo.TVDBID),
		})
	}

	data := TVSeriesMetaData{
		Title:         mediaInfo.TMDBInfo.TVInfo.SerieInfo.Name,         // 电视剧标题
		OriginalTitle: mediaInfo.TMDBInfo.TVInfo.SerieInfo.OriginalName, // 原始标题
		Plot:          mediaInfo.TMDBInfo.TVInfo.SerieInfo.Overview,     // 剧情简介
		Outline:       mediaInfo.TMDBInfo.TVInfo.SerieInfo.Overview,     // 大纲简介
		Actors:        actors,                                           // 演员列表
		Genres:        genres,                                           // 电视剧类型
		Year:          year,                                             // 发行年份
		Premiered:     mediaInfo.TMDBInfo.TVInfo.SerieInfo.FirstAirDate, // 首映日期
		Rating:        mediaInfo.TMDBInfo.TVInfo.SerieInfo.VoteAverage,  // 评分
		Season:        -1,                                               // 默认值 -1，表示未指定季数
		Episode:       -1,                                               // 默认值 -1，表示未指定集数
		TMDBID:        mediaInfo.TMDBID,                                 // TMDB ID
		UniqueIDs:     uniqueIDs,                                        // 唯一ID列表
	}

	return &data
}

func genTVSeasonMetaInfo(mediaInfo *schemas.MediaInfo) *TVSeasonMetaData {
	year, err := strconv.Atoi(mediaInfo.TMDBInfo.TVInfo.SeasonInfo.AirDate[:4])
	if err != nil {
		logrus.Warningf("获取电视剧「第 %d 季 %s」发行年份失败: %v", mediaInfo.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, mediaInfo.TMDBInfo.TVInfo.SeasonInfo.Name, err)
	}

	var uniqueIDs []UniqueID
	if mediaInfo.TMDBID != 0 {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "tmdbid",
			Value: strconv.Itoa(mediaInfo.TMDBID),
		})
	}
	if mediaInfo.IMDBID != "" {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "imdbid",
			Value: mediaInfo.IMDBID,
		})
	}
	if mediaInfo.TVDBID != 0 {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "tvdbid",
			Value: strconv.Itoa(mediaInfo.TVDBID),
		})
	}

	data := TVSeasonMetaData{
		Title:        mediaInfo.TMDBInfo.TVInfo.SeasonInfo.Name,         // 季名称
		Plot:         mediaInfo.TMDBInfo.TVInfo.SeasonInfo.Overview,     // 剧情简介
		Outline:      mediaInfo.TMDBInfo.TVInfo.SeasonInfo.Overview,     // 大纲简介
		Year:         year,                                              // 发行年份
		Premiered:    mediaInfo.TMDBInfo.TVInfo.SeasonInfo.AirDate,      // 首映日期
		ReleaseDate:  mediaInfo.TMDBInfo.TVInfo.SeasonInfo.AirDate,      // 发行日期
		SeasonNumber: mediaInfo.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, // 季数
		TMDBID:       mediaInfo.TMDBID,                                  // TMDB ID
		UniqueIDs:    uniqueIDs,                                         // 唯一ID列表
	}
	return &data
}

func genTVEpisodeMetaInfo(mediaInfo *schemas.MediaInfo) *TVEpisodeMetaData {
	year, err := strconv.Atoi(mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.AirDate[:4])
	if err != nil {
		logrus.Warningf("获取电视剧集「第 %d 季 第 %d 集 %s」发行年份失败: %v", mediaInfo.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.EpisodeNumber, mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.Name, err)
	}
	var (
		actors    []Actor   // 演员列表
		directors []Creator // 导演列表
		writers   []Creator // 编剧列表
		credits   []Creator // 其他制作人员列表
	)

	for _, guestStar := range mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.GuestStars {
		actors = append(actors, Actor{
			Name:    guestStar.Name,
			Role:    guestStar.Character,
			Type:    "GuestStar",
			TMDBID:  strconv.Itoa(guestStar.ID),
			Thumb:   tmdb_controller.GetImageURL(guestStar.ProfilePath),
			Profile: fmt.Sprintf("https://www.themoviedb.org/person/%d", guestStar.ID),
		})
	}

	for _, crew := range mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.Crew {
		creator := Creator{
			Name:    crew.Name,
			TMDBID:  strconv.Itoa(crew.ID),
			Profile: fmt.Sprintf("https://www.themoviedb.org/person/%d", crew.ID),
		}
		switch crew.Job {
		case "Director":
			directors = append(directors, creator)
		case "Writer":
			writers = append(writers, creator)
		default:
			credits = append(credits, creator)
		}
	}

	var uniqueIDs []UniqueID
	if mediaInfo.TMDBID != 0 {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "tmdbid",
			Value: strconv.Itoa(mediaInfo.TMDBID),
		})
	}
	if mediaInfo.IMDBID != "" {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "imdbid",
			Value: mediaInfo.IMDBID,
		})
	}
	if mediaInfo.TVDBID != 0 {
		uniqueIDs = append(uniqueIDs, UniqueID{
			Type:  "tvdbid",
			Value: strconv.Itoa(mediaInfo.TVDBID),
		})
	}

	data := TVEpisodeMetaData{
		Title:     mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.Name,          // 集名称
		Plot:      mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.Overview,      // 剧情简介
		Outline:   mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.Overview,      // 大纲简介
		Year:      year,                                                // 发行年份
		Actors:    actors,                                              // 演员列表
		Directors: directors,                                           // 导演列表
		Writers:   writers,                                             // 编剧列表
		Credits:   credits,                                             // 其他制作人员列表
		Rating:    mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.VoteAverage,   // 评分
		Season:    mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.SeasonNumber,  // 季数
		Episode:   mediaInfo.TMDBInfo.TVInfo.EpisodeInfo.EpisodeNumber, // 集数
		TMDBID:    mediaInfo.TMDBID,                                    // TMDB ID
		UniqueIDs: uniqueIDs,                                           // 唯一ID列表
	}
	return &data
}
func GenMetaDataNFO(infoType InfoType, mediaInfo *schemas.MediaInfo) ([]byte, error) {
	var nfoData InfoData
	switch infoType {
	case InfoTypeMovie:
		if mediaInfo.MediaType != meta.MediaTypeMovie || mediaInfo.TMDBInfo.MovieInfo == nil {
			return nil, fmt.Errorf("媒体信息不完整，无法生成电影 NFO")
		}
		nfoData = genMovieMetaInfo(mediaInfo)
	case InfoTypeTV:
		if mediaInfo.MediaType != meta.MediaTypeTV || mediaInfo.TMDBInfo.TVInfo.SerieInfo == nil {
			return nil, fmt.Errorf("媒体信息不完整，无法生成电视剧 NFO")
		}
		nfoData = genTVSerieMetaInfo(mediaInfo)
	case InfoTypeTVSeason:
		if mediaInfo.MediaType != meta.MediaTypeTV ||
			mediaInfo.TMDBInfo.TVInfo.SeasonInfo == nil ||
			mediaInfo.TMDBInfo.TVInfo.SeasonNumber == -1 {
			return nil, fmt.Errorf("媒体信息不完整，无法生成电视剧季集 NFO")
		}
		nfoData = genTVSeasonMetaInfo(mediaInfo)
	case InfoTypeTVEpisode:
		if mediaInfo.MediaType != meta.MediaTypeTV ||
			mediaInfo.TMDBInfo.TVInfo.EpisodeInfo == nil ||
			mediaInfo.TMDBInfo.TVInfo.EpisodeNumber == -1 {
			return nil, fmt.Errorf("媒体信息不完整，无法生成电视剧集 NFO")
		}
		nfoData = genTVEpisodeMetaInfo(mediaInfo)
	default:
		return nil, fmt.Errorf("不支持的媒体类型: %s", mediaInfo.MediaType)
	}
	return nfoData.XML()
}
