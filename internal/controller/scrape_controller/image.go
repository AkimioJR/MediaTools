package scrape_controller

import (
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func Scrape(dstPath string, info *schemas.MediaInfo) error {
	switch info.MediaType {
	case meta.MediaTypeMovie:
		ScrapeMovieInfo(dstPath, info)
		ScrapeMovieImage(dstPath, info)

	case meta.MediaTypeTV:
		ScrapeTVInfo(dstPath, info)
		ScrapeTVImage(dstPath, info)

	default:
		return fmt.Errorf("不支持的媒体类型: %s", info.MediaType)
	}

	return nil
}

func ScrapeMovieInfo(dstPath string, info *schemas.MediaInfo) {
	metaData := genMovieMetaInfo(info)
	xmlData, err := metaData.XML()
	if err != nil {
		logrus.Errorf("生成电影「%s」元数据 XML 失败: %v", info.TMDBInfo.MovieInfo.Title, err)
	} else {
		infoPath := filepath.Join(filepath.Dir(dstPath), "movie.info")
		file, err := os.Create(infoPath)
		if err != nil {
			logrus.Errorf("创建电影「%s」元数据文件失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		} else {
			defer file.Close()
			_, err = file.Write(xmlData)
			if err != nil {
				logrus.Errorf("写入电影「%s」元数据文件失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			}
		}
	}
}

func ScrapeMovieImage(dstPath string, info *schemas.MediaInfo) {
	err := DownloadTMDBImageAndSave(info.TMDBInfo.MovieInfo.PosterPath, ChangeExt(dstPath, "")+"-poster")
	if err != nil {
		logrus.Errorf("刮削电影「%s」海报失败: %v", info.TMDBInfo.MovieInfo.Title, err)
	}

	movieImage, err := tmdb_controller.GetMovieImage(info.TMDBID)
	if err != nil {
		logrus.Errorf("获取电影「%s」图片信息失败: %v", info.TMDBInfo.MovieInfo.Title, err)
	}

	if len(movieImage.Backdrops) > 0 { // 剧照
		err = DownloadTMDBImageAndSave(movieImage.Backdrops[0].FilePath, filepath.Join(filepath.Dir(dstPath), "backdrop"))
		if err != nil {
			logrus.Errorf("刮削电影「%s」剧照失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}
	}

	if len(movieImage.Posters) > 0 { // 海报
		err = DownloadTMDBImageAndSave(movieImage.Posters[0].FilePath, filepath.Join(filepath.Dir(dstPath), "poster"))
		if err != nil {
			logrus.Errorf("刮削电影「%s」海报失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}
	}

	if len(movieImage.Logos) > 0 { // Logo
		err = DownloadTMDBImageAndSave(movieImage.Logos[0].FilePath, filepath.Join(filepath.Dir(dstPath), "logo"))
		if err != nil {
			logrus.Errorf("刮削电影「%s」Logo 失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}
	}

	{ // 刮削 Fanart 图片
		fanartImagesData, err := fanart_controller.GetMovieImagesData(info.IMDBID)
		if err != nil {
			logrus.Errorf("获取电影「%s」Fanart 图片信息失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}

		if len(fanartImagesData.MovieBackground) > 0 { // 背景图
			err = DownloadFanartImageAndSave(fanartImagesData.MovieBackground[0].URL, filepath.Join(filepath.Dir(dstPath), "background"))
			if err != nil {
				logrus.Errorf("刮削电影「%s」Fanart 背景图失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			}
		}

		if len(fanartImagesData.MovieBanner) > 0 { // 横幅
			err = DownloadFanartImageAndSave(fanartImagesData.MovieBanner[0].URL, filepath.Join(filepath.Dir(dstPath), "banner"))
			if err != nil {
				logrus.Errorf("刮削电影「%s」Fanart 横幅图失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			}
		}

		fanartClearArtPath := filepath.Join(filepath.Dir(dstPath), "clearart")
		if len(fanartImagesData.HDMovieClearArt) > 0 { // Clear Art
			err = DownloadFanartImageAndSave(fanartImagesData.HDMovieClearArt[0].URL, fanartClearArtPath)
		} else if len(fanartImagesData.MovieArt) > 0 { // Clear Art
			err = DownloadFanartImageAndSave(fanartImagesData.MovieArt[0].URL, fanartClearArtPath)
		}
		if err != nil {
			logrus.Errorf("刮削电影「%s」Fanart Clear Art 图片失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}

		if len(fanartImagesData.MovieDisc) > 0 { // 光盘
			err = DownloadFanartImageAndSave(fanartImagesData.MovieDisc[0].URL, filepath.Join(filepath.Dir(dstPath), "disc"))
			if err != nil {
				logrus.Errorf("刮削电影「%s」Fanart 光盘图片失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			}
		}

		if len(fanartImagesData.MovieThumb) > 0 { // 缩略图
			err = DownloadFanartImageAndSave(fanartImagesData.MovieThumb[0].URL, filepath.Join(filepath.Dir(dstPath), "thumb"))
			if err != nil {
				logrus.Errorf("刮削电影「%s」Fanart 缩略图失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			}
		}
	}
}

func ScrapeTVInfo(dstPath string, info *schemas.MediaInfo) {
	tvSeasonDir := filepath.Dir(dstPath)
	tvSerieDir := filepath.Dir(tvSeasonDir)

	serieMetaData := genTVSerieMetaInfo(info)
	xmlData, err := serieMetaData.XML()
	if err != nil {
		logrus.Errorf("生成电视剧「%s」元数据 XML 失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
	} else {
		infoPath := filepath.Join(tvSerieDir, "tv.info")
		file, err := CreateFile(infoPath)
		if err != nil {
			logrus.Errorf("创建电视剧「%s」元数据文件失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
		} else {
			defer file.Close()
			_, err = file.Write(xmlData)
			if err != nil {
				logrus.Errorf("写入电视剧「%s」元数据文件失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
			}
		}
	}

	seasonMetaData := genTVSeasonMetaInfo(info)
	xmlData, err = seasonMetaData.XML()
	if err != nil {
		logrus.Errorf("生成电视剧「%s」第 %d 季元数据 XML 失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, err)
	} else {
		infoPath := filepath.Join(tvSeasonDir, "season.info")
		file, err := CreateFile(infoPath)
		if err != nil {
			logrus.Errorf("创建电视剧「%s」第 %d 季元数据文件失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, err)
		} else {
			defer file.Close()
			_, err = file.Write(xmlData)
			if err != nil {
				logrus.Errorf("写入电视剧「%s」第 %d 季元数据文件失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, err)
			}
		}
	}

	episodeMetaData := genTVEpisodeMetaInfo(info)
	xmlData, err = episodeMetaData.XML()
	if err != nil {
		logrus.Errorf("生成电视剧「%s」第 %d 季第 %d 集元数据 XML 失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, info.TMDBInfo.TVInfo.EpisodeInfo.EpisodeNumber, err)
	} else {
		infoPath := ChangeExt(dstPath, ".info")
		file, err := CreateFile(infoPath)
		if err != nil {
			logrus.Errorf("创建电视剧「%s」第 %d 季第 %d 集元数据文件失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, info.TMDBInfo.TVInfo.EpisodeInfo.EpisodeNumber, err)
		} else {
			defer file.Close()
			_, err = file.Write(xmlData)
			if err != nil {
				logrus.Errorf("写入电视剧「%s」第 %d 季第 %d 集元数据文件失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, info.TMDBInfo.TVInfo.EpisodeInfo.EpisodeNumber, err)
			}
		}
	}
}

func ScrapeTVImage(dstPath string, info *schemas.MediaInfo) {
	tvSeasonDir := filepath.Dir(dstPath)
	tvSerieDir := filepath.Dir(tvSeasonDir)
	{ // 集照片
		err := DownloadTMDBImageAndSave(info.TMDBInfo.TVInfo.EpisodeInfo.StillPath, ChangeExt(dstPath, ""))
		if err != nil {
			logrus.Errorf("刮削电视剧「%s」第 %d 季第 %d 集剧照失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, info.TMDBInfo.TVInfo.EpisodeInfo.EpisodeNumber, err)
		}
	}

	{ // 季照片
		var seasonPosterName string
		switch info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber {
		case 0: // 特别篇
			seasonPosterName = "special-specials-poster"
		default:
			seasonPosterName = fmt.Sprintf("season%02d-poster", info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber)
		}
		seasonPosterPath := filepath.Join(tvSerieDir, seasonPosterName)
		err := DownloadTMDBImageAndSave(info.TMDBInfo.TVInfo.SeasonInfo.PosterPath, seasonPosterPath)
		if err != nil {
			logrus.Errorf("刮削电视剧「%s」第 %d 季海报失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, err)
		}
	}

	err := DownloadTMDBImageAndSave(info.TMDBInfo.TVInfo.SerieInfo.BackdropPath, filepath.Join(tvSerieDir, "backdrop"))
	if err != nil {
		logrus.Errorf("刮削电视剧「%s」剧照失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
	}

	{ // 其他 TMDB 图片
		serieImages, err := tmdb_controller.GetTVSerieImage(info.TMDBID)
		if err != nil {
			logrus.Errorf("获取电视剧「%s」图片信息失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
		} else {
			if len(serieImages.Posters) > 0 { // 海报
				err = DownloadTMDBImageAndSave(serieImages.Posters[0].FilePath, filepath.Join(tvSerieDir, "poster"))
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」海报失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}

			if len(serieImages.Logos) > 0 { // Logo
				err = DownloadTMDBImageAndSave(serieImages.Logos[0].FilePath, filepath.Join(tvSerieDir, "logo"))
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」Logo 失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}
		}
	}

	{ // 刮削 Fanart 图片
		fanartImagesData, err := fanart_controller.GetTVImagesData(info.TVDBID)
		if err != nil {
			logrus.Errorf("获取电视剧「%s」Fanart 图片信息失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
		} else {
			if len(fanartImagesData.ShowBackground) > 0 { // 背景图
				err = DownloadFanartImageAndSave(fanartImagesData.ShowBackground[0].URL, filepath.Join(tvSerieDir, "background"))
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」Fanart 背景图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}

			if len(fanartImagesData.TVBanner) > 0 { // 横幅
				err = DownloadFanartImageAndSave(fanartImagesData.TVBanner[0].URL, filepath.Join(tvSerieDir, "banner"))
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」Fanart 横幅图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}

			if len(fanartImagesData.CharacterArt) > 0 { // 角色图
				err = DownloadFanartImageAndSave(fanartImagesData.CharacterArt[0].URL, filepath.Join(tvSerieDir, "characterart"))
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」Fanart 角色图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}

			cleanArtPath := filepath.Join(tvSerieDir, "clearart")
			if len(fanartImagesData.HDClearArt) > 0 { // HD clearart
				err = DownloadFanartImageAndSave(fanartImagesData.HDClearArt[0].URL, cleanArtPath)
			} else if len(fanartImagesData.ClearArt) > 0 { // clearart
				err = DownloadFanartImageAndSave(fanartImagesData.ClearArt[0].URL, cleanArtPath)
			}
			if err != nil {
				logrus.Errorf("刮削电视剧「%s」Fanart 清晰艺术图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
			}

			if len(fanartImagesData.TVThumb) > 0 { // 缩略图
				err = DownloadFanartImageAndSave(fanartImagesData.TVThumb[0].URL, filepath.Join(tvSerieDir, "thumb"))
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」Fanart 缩略图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}
		}
	}
}
