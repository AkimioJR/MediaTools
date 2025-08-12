package tmdb_controller

import (
	"MediaTools/internal/pkg/themoviedb/v3"
	"image"

	"github.com/sirupsen/logrus"
)

func GetImageURL(path string) string {
	lock.RLock()
	defer lock.RUnlock()

	return client.GetImageURL(path)
}

func DownloadImage(path string) (image.Image, error) {
	lock.RLock()
	defer lock.RUnlock()

	return client.DownloadImage(path)
}

func GetMovieImage(movieID int) (*themoviedb.MovieImage, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电影（TMDB ID: %d）图片", movieID)
	return client.GetMovieImage(movieID, nil, nil)
}

func GetTVSerieImage(tvID int) (*themoviedb.TVSerieImage, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）图片", tvID)
	return client.GetTVSerieImage(tvID, nil, nil)
}

func GetTVSeasonImage(tvID, seasonNumber int) (*themoviedb.TVSeasonImage, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）S%02d图片", tvID, seasonNumber)
	return client.GetTVSeasonImage(tvID, seasonNumber, nil)
}

func GetTVEpisodeImage(tvID, seasonNumber, episodeNumber int) (*themoviedb.TVEpisodeImage, error) {
	lock.RLock()
	defer lock.RUnlock()

	logrus.Infof("开始获取电视剧（TMDB ID: %d）S%02dE%02d图片", tvID, seasonNumber, episodeNumber)
	return client.GetTVEpisodeImage(tvID, seasonNumber, episodeNumber, nil)
}
