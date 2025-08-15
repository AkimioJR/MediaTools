package scrape_controller

import (
	"MediaTools/internal/controller/fanart_controller"
	"MediaTools/internal/controller/recognize_controller"
	"MediaTools/internal/schemas/storage"

	"MediaTools/internal/controller/storage_controller"
	"MediaTools/internal/controller/tmdb_controller"
	"MediaTools/internal/pkg/meta"
	"MediaTools/internal/schemas"
	"MediaTools/utils"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/sirupsen/logrus"
)

func Scrape(dstFile *storage.FileInfo, info *schemas.MediaInfo) error {
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
func RecognizeAndScrape(dstFile *storage.FileInfo, mediaType meta.MediaType, tmdbID int) error {
	videoMeta, _, _ := recognize_controller.ParseVideoMeta(dstFile.Name)
	if mediaType != meta.MediaTypeUnknown && videoMeta.MediaType == meta.MediaTypeUnknown {
		videoMeta.MediaType = mediaType
	}
	if tmdbID != 0 && videoMeta.TMDBID == 0 {
		videoMeta.TMDBID = tmdbID
	}

	info, err := tmdb_controller.RecognizeAndEnrichMedia(videoMeta)
	if err != nil {
		return fmt.Errorf("识别媒体信息失败: %v", err)
	}
	return Scrape(dstFile, info)
}

func ScrapeMovieInfo(dstFile *storage.FileInfo, info *schemas.MediaInfo) {
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

func ScrapeMovieImage(dstFile *storage.FileInfo, info *schemas.MediaInfo) {
	var (
		wg    sync.WaitGroup
		errCh = make(chan error, 10)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := DownloadTMDBImageAndSave(info.TMDBInfo.MovieInfo.PosterPath, utils.ChangeExt(dstFile.Path, "")+"-poster", dstFile.StorageType)
		if err != nil {
			errCh <- fmt.Errorf("刮削电影「%s」海报失败: %v", info.TMDBInfo.MovieInfo.Title, err)
		}
	}()

	wg.Add(1)
	go func() { // 刮削 TMDB 图片
		defer wg.Done()
		movieImage, err := tmdb_controller.GetMovieImage(info.TMDBID)
		if err != nil {
			errCh <- fmt.Errorf("获取电影「%s」图片信息失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			return
		}

		wg.Add(3)
		go func() {
			defer wg.Done()
			if len(movieImage.Backdrops) > 0 { // 剧照
				var paths []string
				for _, backdrop := range movieImage.Backdrops {
					paths = append(paths, backdrop.FilePath)
				}
				path := getSupportImage(paths)
				if path == "" {
					errCh <- fmt.Errorf("电影「%s」没有合适的剧照", info.TMDBInfo.MovieInfo.Title)
				} else {
					err = DownloadTMDBImageAndSave(path, filepath.Join(filepath.Dir(dstFile.Path), "backdrop"), dstFile.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电影「%s」剧照失败: %v", info.TMDBInfo.MovieInfo.Title, err)
					}
				}
			}
		}()

		go func() {
			defer wg.Done()
			if len(movieImage.Posters) > 0 { // 海报
				var paths []string
				for _, poster := range movieImage.Posters {
					paths = append(paths, poster.FilePath)
				}
				path := getSupportImage(paths)
				if path == "" {
					errCh <- fmt.Errorf("电影「%s」没有合适的海报", info.TMDBInfo.MovieInfo.Title)
				} else {
					err = DownloadTMDBImageAndSave(path, filepath.Join(filepath.Dir(dstFile.Path), "poster"), dstFile.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电影「%s」海报失败: %v", info.TMDBInfo.MovieInfo.Title, err)
					}
				}
			}
		}()

		go func() {
			defer wg.Done()
			if len(movieImage.Logos) > 0 { // Logo
				var paths []string
				for _, logo := range movieImage.Logos {
					paths = append(paths, logo.FilePath)
				}
				path := getSupportImage(paths)
				if path == "" {
					errCh <- fmt.Errorf("电影「%s」没有合适的 Logo", info.TMDBInfo.MovieInfo.Title)
				} else {
					err = DownloadTMDBImageAndSave(path, filepath.Join(filepath.Dir(dstFile.Path), "logo"), dstFile.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电影「%s」Logo 失败: %v", info.TMDBInfo.MovieInfo.Title, err)
					}
				}
			}
		}()
	}()

	wg.Add(1)
	go func() { // 刮削 Fanart 图片
		defer wg.Done()
		fanartImagesData, err := fanart_controller.GetMovieImagesData(info.IMDBID)
		if err != nil {
			errCh <- fmt.Errorf("获取电影「%s」Fanart 图片信息失败: %v", info.TMDBInfo.MovieInfo.Title, err)
			return
		}

		wg.Add(5)
		go func() {
			defer wg.Done()
			if len(fanartImagesData.MovieBackground) > 0 { // 背景图
				var urls []string
				for _, bg := range fanartImagesData.MovieBackground {
					urls = append(urls, bg.URL)
				}
				url := getSupportImage(urls)
				if url == "" {
					errCh <- fmt.Errorf("电影「%s」没有合适的 Fanart 背景图", info.TMDBInfo.MovieInfo.Title)
				} else {
					err = DownloadFanartImageAndSave(url, filepath.Join(filepath.Dir(dstFile.Path), "background"), dstFile.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电影「%s」Fanart 背景图失败: %v", info.TMDBInfo.MovieInfo.Title, err)
					}
				}
			}
		}()

		go func() {
			defer wg.Done()
			if len(fanartImagesData.MovieBanner) > 0 { // 横幅
				var urls []string
				for _, banner := range fanartImagesData.MovieBanner {
					urls = append(urls, banner.URL)
				}
				url := getSupportImage(urls)
				if url == "" {
					errCh <- fmt.Errorf("电影「%s」没有合适的 Fanart 横幅图", info.TMDBInfo.MovieInfo.Title)
				} else {
					err = DownloadFanartImageAndSave(url, filepath.Join(filepath.Dir(dstFile.Path), "banner"), dstFile.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电影「%s」Fanart 横幅图失败: %v", info.TMDBInfo.MovieInfo.Title, err)
					}
				}
			}
		}()

		go func() {
			defer wg.Done()
			fanartClearArtPath := filepath.Join(filepath.Dir(dstFile.Path), "clearart")
			var urls []string
			if len(fanartImagesData.HDMovieClearArt) > 0 { // Clear Art
				for _, clearArt := range fanartImagesData.HDMovieClearArt {
					urls = append(urls, clearArt.URL)
				}
			} else if len(fanartImagesData.MovieArt) > 0 { // Clear Art
				for _, clearArt := range fanartImagesData.MovieArt {
					urls = append(urls, clearArt.URL)
				}
			}
			url := getSupportImage(urls)
			if url == "" {
				errCh <- fmt.Errorf("电影「%s」没有合适的 Fanart Clear Art 图", info.TMDBInfo.MovieInfo.Title)
			} else {
				err = DownloadFanartImageAndSave(url, fanartClearArtPath, dstFile.StorageType)
				if err != nil {
					errCh <- fmt.Errorf("刮削电影「%s」Fanart Clear Art 图片失败: %v", info.TMDBInfo.MovieInfo.Title, err)
				}
			}
		}()

		go func() {
			defer wg.Done()
			if len(fanartImagesData.MovieDisc) > 0 { // 光盘
				var urls []string
				for _, disc := range fanartImagesData.MovieDisc {
					urls = append(urls, disc.URL)
				}
				url := getSupportImage(urls)
				if url == "" {
					errCh <- fmt.Errorf("电影「%s」没有合适的 Fanart 光盘图", info.TMDBInfo.MovieInfo.Title)
				} else {

					err = DownloadFanartImageAndSave(url, filepath.Join(filepath.Dir(dstFile.Path), "disc"), dstFile.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电影「%s」Fanart 光盘图片失败: %v", info.TMDBInfo.MovieInfo.Title, err)
					}
				}
			}
		}()

		go func() {
			defer wg.Done()
			if len(fanartImagesData.MovieThumb) > 0 { // 缩略图
				var urls []string
				for _, thumb := range fanartImagesData.MovieThumb {
					urls = append(urls, thumb.URL)
				}
				url := getSupportImage(urls)
				if url == "" {
					errCh <- fmt.Errorf("电影「%s」没有合适的 Fanart 缩略图", info.TMDBInfo.MovieInfo.Title)
				} else {
					err = DownloadFanartImageAndSave(url, filepath.Join(filepath.Dir(dstFile.Path), "thumb"), dstFile.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电影「%s」Fanart 缩略图失败: %v", info.TMDBInfo.MovieInfo.Title, err)
					}
				}
			}
		}()
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		logrus.Warning(err)
	}
}

func ScrapeTVInfo(dstFile *storage.FileInfo, info *schemas.MediaInfo) {
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

func ScrapeTVImage(dstFile *storage.FileInfo, info *schemas.MediaInfo) {
	tvSeasonDir := storage_controller.GetParent(dstFile)
	tvSerieDir := storage_controller.GetParent(tvSeasonDir)

	var (
		errCh = make(chan error, 10)
		wg    sync.WaitGroup
	)

	wg.Add(1)
	go func() { // 集照片
		defer wg.Done()
		err := DownloadTMDBImageAndSave(info.TMDBInfo.TVInfo.EpisodeInfo.StillPath, utils.ChangeExt(dstFile.Path, ""), dstFile.StorageType)
		if err != nil {
			errCh <- fmt.Errorf("刮削电视剧「%s」第 %d 季第 %d 集剧照失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, info.TMDBInfo.TVInfo.EpisodeInfo.EpisodeNumber, err)
		}
	}()

	wg.Add(1)
	go func() { // 季照片
		defer wg.Done()
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
			errCh <- fmt.Errorf("刮削电视剧「%s」第 %d 季海报失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, info.TMDBInfo.TVInfo.SeasonInfo.SeasonNumber, err)
		}
	}()

	wg.Add(1)
	go func() { // 电视剧剧照
		defer wg.Done()
		err := DownloadTMDBImageAndSave(info.TMDBInfo.TVInfo.SerieInfo.BackdropPath, storage_controller.Join(tvSerieDir, "backdrop").Path, tvSerieDir.StorageType)
		if err != nil {
			errCh <- fmt.Errorf("刮削电视剧「%s」剧照失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
		}
	}()

	wg.Add(1)
	go func() { // 其他 TMDB 图片
		defer wg.Done()

		serieImages, err := tmdb_controller.GetTVSerieImage(info.TMDBID)
		if err != nil {
			errCh <- fmt.Errorf("获取电视剧「%s」图片信息失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
			return
		}

		wg.Add(2)

		go func() { // 海报
			defer wg.Done()
			if len(serieImages.Posters) > 0 {
				var paths []string
				for _, poster := range serieImages.Posters {
					paths = append(paths, poster.FilePath)
				}
				path := getSupportImage(paths)
				if path == "" {
					errCh <- fmt.Errorf("电视剧「%s」没有合适的海报", info.TMDBInfo.TVInfo.SerieInfo.Name)
				} else {
					err = DownloadTMDBImageAndSave(path, filepath.Join(tvSerieDir.Path, "poster"), tvSerieDir.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电视剧「%s」海报失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
					}
				}
			}
		}()

		go func() { // Logo
			defer wg.Done()
			if len(serieImages.Logos) > 0 {
				var paths []string
				for _, logo := range serieImages.Logos {
					paths = append(paths, logo.FilePath)
				}
				err = DownloadTMDBImageAndSave(getSupportImage(paths), filepath.Join(tvSerieDir.Path, "logo"), tvSerieDir.StorageType)
				if err != nil {
					errCh <- fmt.Errorf("刮削电视剧「%s」Logo 失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}
		}()
	}()

	wg.Add(1)
	go func() { // 刮削 Fanart 图片
		defer wg.Done()

		fanartImagesData, err := fanart_controller.GetTVImagesData(info.TVDBID)
		if err != nil {
			errCh <- fmt.Errorf("获取电视剧「%s」Fanart 图片信息失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
			return
		}

		wg.Add(5)
		go func() { // 背景图
			defer wg.Done()
			if len(fanartImagesData.ShowBackground) > 0 {
				var urls []string
				for _, bg := range fanartImagesData.ShowBackground {
					urls = append(urls, bg.URL)
				}
				url := getSupportImage(urls)
				if url == "" {
					errCh <- fmt.Errorf("电视剧「%s」没有合适的 Fanart 背景图", info.TMDBInfo.TVInfo.SerieInfo.Name)
				} else {
					err = DownloadFanartImageAndSave(url, filepath.Join(tvSerieDir.Path, "background"), tvSerieDir.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电视剧「%s」Fanart 背景图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
					}
				}
			}
		}()

		go func() { // 横幅
			defer wg.Done()
			if len(fanartImagesData.TVBanner) > 0 {
				var urls []string
				for _, banner := range fanartImagesData.TVBanner {
					urls = append(urls, banner.URL)
				}
				url := getSupportImage(urls)
				if url == "" {
					errCh <- fmt.Errorf("电视剧「%s」没有合适的 Fanart 横幅图", info.TMDBInfo.TVInfo.SerieInfo.Name)
				} else {
					err = DownloadFanartImageAndSave(url, filepath.Join(tvSerieDir.Path, "banner"), tvSerieDir.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电视剧「%s」Fanart 横幅图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
					}
				}
			}
		}()

		go func() { // 角色图
			defer wg.Done()
			if len(fanartImagesData.CharacterArt) > 0 {
				var urls []string
				for _, charArt := range fanartImagesData.CharacterArt {
					urls = append(urls, charArt.URL)
				}
				url := getSupportImage(urls)
				if url == "" {
					errCh <- fmt.Errorf("电视剧「%s」没有合适的 Fanart 角色图", info.TMDBInfo.TVInfo.SerieInfo.Name)
				} else {
					err = DownloadFanartImageAndSave(url, filepath.Join(tvSerieDir.Path, "characterart"), tvSerieDir.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电视剧「%s」Fanart 角色图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
					}
				}
			}
		}()

		go func() { // 清晰艺术图
			defer wg.Done()
			cleanArtPath := filepath.Join(tvSerieDir.Path, "clearart")
			var urls []string
			if len(fanartImagesData.HDClearArt) > 0 { // HD clearart
				for _, clearArt := range fanartImagesData.HDClearArt {
					urls = append(urls, clearArt.URL)
				}
			} else if len(fanartImagesData.ClearArt) > 0 { // clearart
				for _, clearArt := range fanartImagesData.HDClearArt {
					urls = append(urls, clearArt.URL)
				}
			}
			url := getSupportImage(urls)
			if url == "" {
				errCh <- fmt.Errorf("电视剧「%s」没有合适的 Fanart Clear Art 图片", info.TMDBInfo.TVInfo.SerieInfo.Name)
			} else {
				err = DownloadFanartImageAndSave(url, cleanArtPath, tvSerieDir.StorageType)
				if err != nil {
					errCh <- fmt.Errorf("刮削电视剧「%s」Fanart 清晰艺术图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
				}
			}
		}()

		go func() { // 缩略图
			defer wg.Done()
			if len(fanartImagesData.TVThumb) > 0 {
				var urls []string
				for _, thumb := range fanartImagesData.TVThumb {
					urls = append(urls, thumb.URL)
				}
				url := getSupportImage(urls)
				if url == "" {
					errCh <- fmt.Errorf("电视剧「%s」没有合适的 Fanart 缩略图", info.TMDBInfo.TVInfo.SerieInfo.Name)
				} else {
					err = DownloadFanartImageAndSave(url, filepath.Join(tvSerieDir.Path, "thumb"), tvSerieDir.StorageType)
					if err != nil {
						errCh <- fmt.Errorf("刮削电视剧「%s」Fanart 缩略图失败: %v", info.TMDBInfo.TVInfo.SerieInfo.Name, err)
					}
				}

			}
		}()
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		logrus.Warning(err)
	}
}
