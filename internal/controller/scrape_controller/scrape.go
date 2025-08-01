package scrape_controller

import (
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"MediaTools/utils"
	"fmt"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func Scrape(dstFile *schemas.FileInfo, info *schemas.MediaInfo) error {
	switch info.MediaType {
	case meta.MediaTypeMovie:
		ScrapeMovieInfo(dstFile, info)
		ScrapeMovieImage(dstFile, info)

	case meta.MediaTypeTV:
		ScrapeTVInfo(dstFile, info)
		ScrapeTVImage(dstFile, info)

	default:
		return fmt.Errorf("不支持的媒体类型: %s", info.MediaType)
	}

	return nil
}

// RecognizeAndScrape 识别并刮削媒体信息
// 识别目标文件的元数据，查询 TMDB 获取媒体信息，并在该文件夹进行刮削
func RecognizeAndScrape(dstFile *schemas.FileInfo, mediaType *meta.MediaType, tmdbID *int) error {
	videoMeta := meta.ParseVideoMeta(dstFile.Name)
	info, err := tmdb_controller.RecognizeAndEnrichMedia(videoMeta, mediaType, nil)
	if err != nil {
		return fmt.Errorf("识别媒体信息失败: %v", err)
	}
	return Scrape(dstFile, info)
}

func ScrapeMovieInfo(dstFile *schemas.FileInfo, info *schemas.MediaInfo) {
	metaData := genMovieMetaInfo(info)
	xmlData, err := metaData.XML()
	if err != nil {
		logrus.Errorf("生成电影「%s」元数据 XML 失败: %v", info.TMDBInfo.MovieInfo.Title, err)
	} else {
		infoFile := storage_controller.Join(storage_controller.GetParent(dstFile), "movie.info")
		reader, err := bytes2Reader(xmlData)
		if err != nil {
			logrus.Warning("创建 reader 失败: %w", err)
			return
		}
		err = storage_controller.CreateFile(infoFile, reader)
		if err != nil {
			logrus.Errorf("创建电影「%s」元数据文件失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}
	}
}

func ScrapeMovieImage(dstFile *schemas.FileInfo, info *schemas.MediaInfo) {
	err := DownloadTMDBImageAndSave(info.TMDBInfo.MovieInfo.PosterPath, utils.ChangeExt(dstFile.Path, "")+"-poster", dstFile.StorageType)
	if err != nil {
		logrus.Errorf("刮削电影「%s」海报失败: %v", info.TMDBInfo.MovieInfo.Title, err)
	}

	movieImage, err := tmdb_controller.GetMovieImage(info.TMDBID)
	if err != nil {
		logrus.Errorf("获取电影「%s」图片信息失败: %v", info.TMDBInfo.MovieInfo.Title, err)
	}

	if len(movieImage.Backdrops) > 0 { // 剧照
		err = DownloadTMDBImageAndSave(movieImage.Backdrops[0].FilePath, filepath.Join(filepath.Dir(dstFile.Path), "backdrop"), dstFile.StorageType)
		if err != nil {
			logrus.Errorf("刮削电影「%s」剧照失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}
	}

	if len(movieImage.Posters) > 0 { // 海报
		err = DownloadTMDBImageAndSave(movieImage.Posters[0].FilePath, filepath.Join(filepath.Dir(dstFile.Path), "poster"), dstFile.StorageType)
		if err != nil {
			logrus.Errorf("刮削电影「%s」海报失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}
	}

	if len(movieImage.Logos) > 0 { // Logo
		err = DownloadTMDBImageAndSave(movieImage.Logos[0].FilePath, filepath.Join(filepath.Dir(dstFile.Path), "logo"), dstFile.StorageType)
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
			err = DownloadFanartImageAndSave(fanartImagesData.MovieBackground[0].URL, filepath.Join(filepath.Dir(dstFile.Path), "background"), dstFile.StorageType)
			if err != nil {
				logrus.Errorf("刮削电影「%s」Fanart 背景图失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			}
		}

		if len(fanartImagesData.MovieBanner) > 0 { // 横幅
			err = DownloadFanartImageAndSave(fanartImagesData.MovieBanner[0].URL, filepath.Join(filepath.Dir(dstFile.Path), "banner"), dstFile.StorageType)
			if err != nil {
				logrus.Errorf("刮削电影「%s」Fanart 横幅图失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			}
		}

		fanartClearArtPath := filepath.Join(filepath.Dir(dstFile.Path), "clearart")
		if len(fanartImagesData.HDMovieClearArt) > 0 { // Clear Art
			err = DownloadFanartImageAndSave(fanartImagesData.HDMovieClearArt[0].URL, fanartClearArtPath, dstFile.StorageType)
		} else if len(fanartImagesData.MovieArt) > 0 { // Clear Art
			err = DownloadFanartImageAndSave(fanartImagesData.MovieArt[0].URL, fanartClearArtPath, dstFile.StorageType)
		}
		if err != nil {
			logrus.Errorf("刮削电影「%s」Fanart Clear Art 图片失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}

		if len(fanartImagesData.MovieDisc) > 0 { // 光盘
			err = DownloadFanartImageAndSave(fanartImagesData.MovieDisc[0].URL, filepath.Join(filepath.Dir(dstFile.Path), "disc"), dstFile.StorageType)
			if err != nil {
				logrus.Errorf("刮削电影「%s」Fanart 光盘图片失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			}
		}

		if len(fanartImagesData.MovieThumb) > 0 { // 缩略图
			err = DownloadFanartImageAndSave(fanartImagesData.MovieThumb[0].URL, filepath.Join(filepath.Dir(dstFile.Path), "thumb"), dstFile.StorageType)
			if err != nil {
				logrus.Errorf("刮削电影「%s」Fanart 缩略图失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			}
		}
	}
}

func ScrapeTVInfo(dstFile *schemas.FileInfo, info *schemas.MediaInfo) {
	tvSeasonDir := storage_controller.GetParent(dstFile)
	tvSerieDir := storage_controller.GetParent(tvSeasonDir)

	serieMetaData := genTVSerieMetaInfo(info)
	xmlData, err := serieMetaData.XML()
	if err != nil {
		logrus.Errorf("生成电视剧「%s」元数据 XML 失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
	} else {
		infoFile := storage_controller.Join(tvSerieDir, "tv.info")
		reader, err := bytes2Reader(xmlData)
		if err != nil {
			logrus.Warningf("创建 reader 失败: %v", err)
		} else {
			err = storage_controller.CreateFile(infoFile, reader)
			if err != nil {
				logrus.Errorf("创建电视剧「%s」元数据文件失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
			}
		}
	}

	seasonMetaData := genTVSeasonMetaInfo(info)
	xmlData, err = seasonMetaData.XML()
	if err != nil {
		logrus.Errorf("生成电视剧「%s」第 %d 季元数据 XML 失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, err)
	} else {
		infoFile := storage_controller.Join(tvSeasonDir, "season.info")
		reader, err := bytes2Reader(xmlData)
		if err != nil {
			logrus.Warning("创建 reader 失败: %w", err)
		} else {
			err = storage_controller.CreateFile(infoFile, reader)
			if err != nil {
				logrus.Errorf("创建电视剧「%s」第 %d 季元数据文件失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, err)
			}
		}
	}

	episodeMetaData := genTVEpisodeMetaInfo(info)
	xmlData, err = episodeMetaData.XML()
	if err != nil {
		logrus.Errorf("生成电视剧「%s」第 %d 季第 %d 集元数据 XML 失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, info.TMDBInfo.TVInfo.EpisodeInfo.EpisodeNumber, err)
	} else {
		infoPath := utils.ChangeExt(dstFile.Path, ".info")
		infoFile, err := storage_controller.GetFile(infoPath, dstFile.StorageType)
		if err != nil {
			logrus.Warningf("获取 %s:/%s 句柄失败:%v", dstFile.StorageType, infoPath, err)
		} else {
			reader, err := bytes2Reader(xmlData)
			if err != nil {
				logrus.Warning("创建 reader 失败: %w", err)
			} else {
				err = storage_controller.CreateFile(infoFile, reader)
				if err != nil {
					logrus.Errorf("创建电视剧「%s」第 %d 季第 %d 集元数据文件失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, info.TMDBInfo.TVInfo.EpisodeInfo.EpisodeNumber, err)
				}
			}
		}
	}
}

func ScrapeTVImage(dstFile *schemas.FileInfo, info *schemas.MediaInfo) {
	tvSeasonDir := storage_controller.GetParent(dstFile)
	tvSerieDir := storage_controller.GetParent(tvSeasonDir)

	{ // 集照片
		err := DownloadTMDBImageAndSave(info.TMDBInfo.TVInfo.EpisodeInfo.StillPath, utils.ChangeExt(dstFile.Path, ""), dstFile.StorageType)
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
		seasonPosterFile := storage_controller.Join(tvSerieDir, seasonPosterName)
		err := DownloadTMDBImageAndSave(info.TMDBInfo.TVInfo.SeasonInfo.PosterPath, seasonPosterFile.Path, seasonPosterFile.StorageType)
		if err != nil {
			logrus.Errorf("刮削电视剧「%s」第 %d 季海报失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, err)
		}
	}

	err := DownloadTMDBImageAndSave(info.TMDBInfo.TVInfo.SerieInfo.BackdropPath, storage_controller.Join(tvSerieDir, "backdrop").Path, tvSerieDir.StorageType)
	if err != nil {
		logrus.Errorf("刮削电视剧「%s」剧照失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
	}

	{ // 其他 TMDB 图片
		serieImages, err := tmdb_controller.GetTVSerieImage(info.TMDBID)
		if err != nil {
			logrus.Errorf("获取电视剧「%s」图片信息失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
		} else {
			if len(serieImages.Posters) > 0 { // 海报
				err = DownloadTMDBImageAndSave(serieImages.Posters[0].FilePath, filepath.Join(tvSerieDir.Path, "poster"), tvSerieDir.StorageType)
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」海报失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}

			if len(serieImages.Logos) > 0 { // Logo
				err = DownloadTMDBImageAndSave(serieImages.Logos[0].FilePath, filepath.Join(tvSerieDir.Path, "logo"), tvSerieDir.StorageType)
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
				err = DownloadFanartImageAndSave(fanartImagesData.ShowBackground[0].URL, filepath.Join(tvSerieDir.Path, "background"), tvSerieDir.StorageType)
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」Fanart 背景图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}

			if len(fanartImagesData.TVBanner) > 0 { // 横幅
				err = DownloadFanartImageAndSave(fanartImagesData.TVBanner[0].URL, filepath.Join(tvSerieDir.Path, "banner"), tvSerieDir.StorageType)
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」Fanart 横幅图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}

			if len(fanartImagesData.CharacterArt) > 0 { // 角色图
				err = DownloadFanartImageAndSave(fanartImagesData.CharacterArt[0].URL, filepath.Join(tvSerieDir.Path, "characterart"), tvSerieDir.StorageType)
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」Fanart 角色图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}

			cleanArtPath := filepath.Join(tvSerieDir.Path, "clearart")
			if len(fanartImagesData.HDClearArt) > 0 { // HD clearart
				err = DownloadFanartImageAndSave(fanartImagesData.HDClearArt[0].URL, cleanArtPath, tvSerieDir.StorageType)
			} else if len(fanartImagesData.ClearArt) > 0 { // clearart
				err = DownloadFanartImageAndSave(fanartImagesData.ClearArt[0].URL, cleanArtPath, tvSerieDir.StorageType)
			}
			if err != nil {
				logrus.Errorf("刮削电视剧「%s」Fanart 清晰艺术图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
			}

			if len(fanartImagesData.TVThumb) > 0 { // 缩略图
				err = DownloadFanartImageAndSave(fanartImagesData.TVThumb[0].URL, filepath.Join(tvSerieDir.Path, "thumb"), tvSerieDir.StorageType)
				if err != nil {
					logrus.Errorf("刮削电视剧「%s」Fanart 缩略图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}
		}
	}
}
